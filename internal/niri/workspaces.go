package niri

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// ListWorkspaces returns a snapshot of all niri workspaces by running
// `niri msg -j workspaces` and unmarshaling its JSON array. It is a one-shot
// call (the daemon uses the event stream for live updates); the `tile-data`
// producer calls it once per invocation to resolve a desktop index to its
// workspace id, output, and focus. The Workspace type is defined in
// eventstream.go and shared with the event-stream model.
func ListWorkspaces() ([]Workspace, error) {
	cmd := exec.Command("niri", "msg", "-j", "workspaces")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("niri msg -j workspaces: %w", err)
	}
	var ws []Workspace
	if err := json.Unmarshal(out, &ws); err != nil {
		return nil, fmt.Errorf("decode niri workspaces: %w", err)
	}
	return ws, nil
}

// MoveWorkspaceToIndex runs `niri msg action move-workspace-to-index --reference
// <src> <dst>`, reordering one desktop; its windows travel with it and niri
// shifts the desktops in between. As with focus-workspace, a bare index
// resolves on the FOCUSED output, so cross with FocusMonitor first to reorder
// another output's desktops.
func MoveWorkspaceToIndex(src, dst int) error {
	return runAction("move-workspace-to-index", "--reference", itoa(src), itoa(dst))
}
