package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io"
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
		SetCurrentNode(root).
		SetGraphics(false)

	var addArray func(node *tview.TreeNode, dec *json.Decoder, delim byte)
	var addNode func(node *tview.TreeNode, dec *json.Decoder, delim byte)

	addArray = func(node *tview.TreeNode, dec *json.Decoder, delim byte) {
		for {
			// Get key or delim
			tok, err := dec.Token()
			if err == io.EOF {
				return
			}

			switch v := tok.(type) {
			case json.Delim:
				newNode := tview.NewTreeNode(v.String())
				node.AddChild(newNode).SetSelectable(true)
				if v.String()[0] == delim {
					return
				} else if v.String() == "{" {
					// create } node
					node.AddChild(tview.NewTreeNode("}")).
						SetSelectable(true)
				} else if v.String() == "[" {
					addNode(newNode, dec, ']')
					node.AddChild(tview.NewTreeNode("]")).
						SetSelectable(true)
				}
			case bool:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", v))).
					SetSelectable(true)
			case float64:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", v))).
					SetSelectable(true)
			case string:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%v", v))).
					SetSelectable(true)
			case nil:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("null"))).
					SetSelectable(true)
			default:
				fmt.Println("Error2: should not be here")
				return
			}
		}
	}

	addNode = func(node *tview.TreeNode, dec *json.Decoder, delim byte) {
		for {

			// Get key or delim
			tok, err := dec.Token()
			if err == io.EOF {
				return
			}

			if _, ok := tok.(json.Delim); ok {
				return // TODO: test if right delim
			}

			// handle error
			keyText, _ := tok.(string)

			// Get value
			tok, err = dec.Token()
			if err == io.EOF {
				break
			}

			// TODO: use real delim type not string
			switch v := tok.(type) {
			case json.Delim:
				newNode := tview.NewTreeNode(v.String())
				node.AddChild(newNode).SetSelectable(true)
				if v.String() == "{" {
					// create } node
					node.AddChild(tview.NewTreeNode("}")).
						SetSelectable(true)
				} else if v.String() == "[" {
					addArray(newNode, dec, ']')
					node.AddChild(tview.NewTreeNode("]")).
						SetSelectable(true)
				}
			case bool:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%s: %v", keyText, v))).
					SetSelectable(true)
			case float64:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%s: %v", keyText, v))).
					SetSelectable(true)
			case string:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%s: %v", keyText, v))).
					SetSelectable(true)
			case nil:
				node.AddChild(tview.NewTreeNode(fmt.Sprintf("%s: null", keyText))).
					SetSelectable(true)
			default:
				fmt.Println("Error2: should not be here")
				return
			}
		}
		return
	}

	addNode(root, json.NewDecoder(file), '}')

	return tree
}

func main() {
	file, _ := os.Open("test.json")
	tree := createJsonTree(file)

	// handle what to do with the file when it is selected
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		//reference := node.GetReference()
		node.SetExpanded(!node.IsExpanded())
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
