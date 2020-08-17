package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
)

func printPrompt() {
	wd, _ := os.Getwd()
	fmt.Printf("⟪%s⟫\nᐉ ", wd)
}

func createJsonTree(file *os.File) *tview.TreeView {
	// Create root
	root := tview.
		NewTreeNode("{").
		SetColor(tcell.ColorRed)

	// Create Tree
	tree := tview.
		NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// Create decoder
	dec := json.NewDecoder(file)

	var addNode func(node *tview.TreeNode, delim byte)
	addNode = func(node *tview.TreeNode, delim byte) {
		for dec.More() {
			// TODO error handle
			tok, _ := dec.Token()

			// TODO: use real delim type not string
			switch v := tok.(type) {
			case json.Delim:
				if v.String()[0] == delim {
					return
				}

				newNode := tview.NewTreeNode(v.String())
				node.AddChild(newNode).SetSelectable(true)

				if v.String() == "{" {
					fmt.Println("beg }")
					addNode(newNode, '}')
					fmt.Println("end }")
				} else if v.String() == "[" {
					fmt.Println("beg ]")
					addNode(newNode, ']')
					fmt.Println("end ]")
				}
			case bool:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", v))).
					SetSelectable(true)
			case float64:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", v))).
					SetSelectable(true)
			case string:
				node.AddChild(tview.NewTreeNode(v)).
					SetSelectable(true)
			case nil:
				node.AddChild(tview.NewTreeNode("null")).
					SetSelectable(true)
			}
		}
	}

	addNode(root, '}')

	return tree
}

func main() {
	file, _ := os.Open("test.json")
	tree := createJsonTree(file)

	// handle what to do with the file when it is selected
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		//reference := node.GetReference()
		node.SetExpanded(!node.IsExpanded())
		fmt.Println("test")
	})

	// handle additional key presses
	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ' ' {
			//	name := tree.GetCurrentNode().GetReference()
			return nil
		}
		return event
	})

	// Create app
	app := tview.
		NewApplication().
		SetRoot(tree, true)

	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		return false
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
