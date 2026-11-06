package tile

// This file implements `claude-status tile-move <src> <dst>` — the command a
// pwetty desktop tile's on_drop runs when one tile is dragged onto another.
// Dropping tile <src> on tile <dst> moves that desktop to index <dst>; niri
// shifts the ones in between (so dropping 2 on 1 swaps them) and the windows
// travel with their desktop. The tiles re-render from the daemon's next
// workspace event, nothing here touches the cache.

import (
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/mrzor/claude-status/internal/niri"
)

// RunMove executes the tile-move subcommand. args is os.Args[2:]. Flags mirror
// RunToggle's, so a bar pins its output the same way its tile-watch does
// (`tile-move --output eDP-1 2 1`; the flag must precede the indexes).
func RunMove(args []string) error {
	output, src, dst, err := parseMove(args)
	if err != nil || src == dst {
		return err
	}
	// A bare index resolves on the focused output (see jump): cross first.
	focusedOutput, _ := desktopFocus(output, src)
	if crossesOutput(output, focusedOutput) {
		_ = niri.FocusMonitor(output)
	}
	return niri.MoveWorkspaceToIndex(src, dst)
}

// parseMove reads `[--output <connector>] <src> <dst>`; both are 1-based
// desktop indexes.
func parseMove(args []string) (output string, src, dst int, err error) {
	output = defaultOutput
	fs := flag.NewFlagSet("tile-move", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&output, "output", output, "niri output (connector) the tiles' desktops live on")
	_ = fs.Parse(args)
	src, errSrc := strconv.Atoi(fs.Arg(0))
	dst, errDst := strconv.Atoi(fs.Arg(1))
	if errSrc != nil || errDst != nil || src < 1 || dst < 1 {
		return "", 0, 0, fmt.Errorf("usage: tile-move [--output <connector>] <src-idx> <dst-idx>")
	}
	return output, src, dst, nil
}
