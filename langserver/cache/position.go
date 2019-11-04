// Copyright 2019 Tobias Guggenmos
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"fmt"
	"go/token"
	"os"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/span"
)

// PositionToProtocolPostion converts a token.Position to a protocol.Position
func (d *Document) PositionToProtocolPostion(ctx context.Context, pos token.Position) (protocol.Position, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case <-ctx.Done():
		return protocol.Position{}, ctx.Err()
	default:
		line := pos.Line
		char := pos.Column

		// Can happen when parsing empty files
		if line < 1 {
			return protocol.Position{
				Line:      0,
				Character: 0,
			}, nil
		}

		// Convert to the Postions as described in the LSP Spec
		// LineStart can panic
		offset := int(d.posData.LineStart(line)) - d.posData.Base() + char - 1
		point := span.NewPoint(line, char, offset)

		var err error

		char, err = span.ToUTF16Column(point, []byte(d.content))
		// Protocol has zero based positions
		char--
		line--

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return protocol.Position{}, err
		}

		return protocol.Position{
			Line:      float64(line),
			Character: float64(char),
		}, nil
	}
}

// ProtocolPositionToTokenPos converts a token.Pos to a protocol.Position
func (d *Document) ProtocolPositionToTokenPos(ctx context.Context, pos protocol.Position) (token.Pos, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		// protocol.Position is 0 based
		line := int(pos.Line) + 1
		char := int(pos.Character)
		offset := int(d.posData.LineStart(line)) - d.posData.Base()
		point := span.NewPoint(line, 1, offset)
		point, err := span.FromUTF16Column(point, char, []byte(d.content))

		if err != nil {
			return token.NoPos, err
		}

		char = point.Column()

		return d.posData.LineStart(line) + token.Pos(char), nil
	}
}

// EndOfLine returns the end of the Line of the given protocol.Position
func EndOfLine(p protocol.Position) protocol.Position {
	return protocol.Position{
		Line:      p.Line + 1,
		Character: 0,
	}
}
