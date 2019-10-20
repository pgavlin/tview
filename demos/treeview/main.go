// Demo code for the TreeView primitive.
package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Show a navigable tree view of the current directory.
func main() {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetBorder(true)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			columns := []*tview.TreeCell{
				tview.NewTreeCell(fmt.Sprintf("%v", file.Size())).SetMaxWidth(20),
				tview.NewTreeCell(file.ModTime().String()).SetMaxWidth(20),
			}

			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir()).
				SetColumns(columns)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
				columns[0].Text = ""
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
