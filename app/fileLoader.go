package app

import (
	"errors"
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"math"
	"os"
)

type FileLoadLayer struct {
	origin Vec2Df32
	size   Vec2Df32
	ltNode *LayerTreeNode

	scrollVal     int32
	workingDir    string
	fileEntries   []os.DirEntry
	selectedEntry int
	entrySize     int
	fileToLoad    string
	callback      func(fname string)
}

func (f *FileLoadLayer) OnCreate() {
	f.scrollVal = 0
	f.workingDir = "./"
	f.selectedEntry = 0
	f.entrySize = 32
	log.Printf("FileLoadLayer Layer created with unique id: %v\n", f.ltNode.UniqueID)
}

func (f *FileLoadLayer) OnRemove() {
	log.Printf("FileLoadLayer Layer removed with unique id: %v\n", f.ltNode.UniqueID)
}

func (f *FileLoadLayer) SetLTNode(ltNode *LayerTreeNode){
	f.ltNode = ltNode
}

func (f *FileLoadLayer) OnEvent() {}
func (f *FileLoadLayer) OnUpdate() {
	f.updateFileEntries()
	f.updateSelectedEntry()
	err := f.processClicks()
	if err != nil {
		log.Fatal(err)
	}
}
func (f *FileLoadLayer) OnRender() {
	f.drawFileDialog()
}
func (f *FileLoadLayer) SetTransform(origin, size Vec2Df32) {
	f.origin = origin
	f.size = size
}
func (f *FileLoadLayer) GetTransform() (Vec2Df32, Vec2Df32){
	return f.origin, f.size
}
func (f *FileLoadLayer) SetCallback(c func(fname string)) {
	f.callback = c
}

func (f *FileLoadLayer) updateFileEntries() {
	//get the current dir entries
	dirEntries, err := os.ReadDir(f.workingDir)
	if err != nil {
		log.Fatal(err)
	}
	f.fileEntries = nil
	for _, entry := range dirEntries {
		f.fileEntries = append(f.fileEntries, entry)
	}
}

func (f *FileLoadLayer) drawFileDialog() {

	frame := f.ltNode.GetFrame()

	//draw the file dialog
	fileDialogRect := rl.Rectangle{(frame.X),
								   (frame.Y),
								   (frame.Width),
		                           (frame.Height)}
	fileDialogColor := rl.NewColor(77, 77, 77, 200)
	rl.DrawRectangle(int32(fileDialogRect.X),
		int32(fileDialogRect.Y),
		int32(fileDialogRect.Width),
		int32(fileDialogRect.Height), fileDialogColor)

	//draw the scroll bar
	var sliderOrigin, sliderSize Vec2Df32
	sliderOrigin.X = frame.X + frame.Width
	sliderOrigin.Y = frame.Y
	sliderSize.X = 0.02 * fileDialogRect.Width
	sliderSize.Y = fileDialogRect.Height

	scrollBarBounds := rl.Rectangle{sliderOrigin.X, sliderOrigin.Y, sliderSize.X, sliderSize.Y}
	f.scrollVal = gui.ScrollBar(scrollBarBounds, f.scrollVal, 0, int32(len(f.fileEntries)+1))

	//draw the selected entry rectangle
	idx := f.selectedEntry
	if idx != -1 {
		selectedEntryColor := rl.NewColor(200, 200, 200, 200)
		selectedEntryRect := rl.Rectangle{fileDialogRect.X,
			fileDialogRect.Y + float32((idx-int(f.scrollVal))*f.entrySize),
			fileDialogRect.Width, float32(f.entrySize)}
		rl.DrawRectangle(int32(selectedEntryRect.X),
			int32(selectedEntryRect.Y),
			int32(selectedEntryRect.Width),
			int32(selectedEntryRect.Height),
			selectedEntryColor)
	}

	//draw the available files
	rem := int(fileDialogRect.Height) % f.entrySize
	minY := int(fileDialogRect.Y)
	maxY := int(math.Min(float64(minY+f.entrySize*(len(f.fileEntries)+1)),
		float64(minY+int(fileDialogRect.Height)-rem)))
	if int(f.scrollVal) == 0 {
		rl.DrawText("(dir) ../",
			int32(fileDialogRect.X),
			int32(fileDialogRect.Y)+int32(f.entrySize*(0-int(f.scrollVal))),
			int32(f.entrySize), rl.White)
	}
	for i, dir_entry := range f.fileEntries {
		draw_idx := i - int(f.scrollVal)
		yval := int(int32(fileDialogRect.Y)+int32(f.entrySize*(draw_idx+1))) + f.entrySize
		if yval > maxY || yval < minY+f.entrySize {
			continue
		}
		var name string = dir_entry.Name()
		if dir_entry.IsDir() {
			name = "(dir) " + name
		}
		rl.DrawText(name,
			int32(fileDialogRect.X),
			int32(fileDialogRect.Y)+int32(f.entrySize*(draw_idx+1)),
			int32(f.entrySize), rl.White)
	}
}

func (f *FileLoadLayer) updateSelectedEntry() {
	var mouseVec rl.Vector2 = rl.GetMousePosition()

	frame := f.ltNode.GetFrame()

	//draw the file dialog
	fileDialogRect := rl.Rectangle{(frame.X),
		(frame.Y),
		(frame.Width),
		(frame.Height)}

	minY := int(fileDialogRect.Y)
	rem := int(fileDialogRect.Height) % f.entrySize
	maxY := int(math.Min(float64(minY+f.entrySize*(len(f.fileEntries)+1)),
		float64(minY+int(fileDialogRect.Height)-rem)))
	actualY := int(mouseVec.Y)
	inDialogY := mouseVec.Y - fileDialogRect.Y

	minX := int(fileDialogRect.X)
	maxX := minX + int(fileDialogRect.Width)
	actualX := int(mouseVec.X)
	if actualX < minX || actualX > maxX {
		f.selectedEntry = -1
		return
	}

	if actualY >= maxY {
		f.selectedEntry = -1
	} else if actualY < minY {
		f.selectedEntry = -1
	} else {
		idx := int(f.scrollVal) + (int(inDialogY) / f.entrySize)
		if idx <= len(f.fileEntries) {
			f.selectedEntry = int(f.scrollVal) + (int(inDialogY) / f.entrySize)
		} else {
			f.selectedEntry = -1
		}
	}
}

func (f *FileLoadLayer) processClicks() error {
	if rl.IsMouseButtonReleased(0) {
		if f.selectedEntry > len(f.fileEntries) {
			return errors.New("selected entry idx is greater than the number of entries")
		}
		if f.selectedEntry == -1 {
			return nil
		}

		if f.selectedEntry > 0 {
			var idx int = f.selectedEntry - 1
			entry := f.fileEntries[idx]
			if entry.IsDir() {
				f.workingDir = f.workingDir + "/" + entry.Name()
				f.scrollVal = 0
			} else {
				//run the callback on the file path and remove this layer
				f.fileToLoad = f.workingDir + "/" + entry.Name()
				f.ltNode.Remove()
				f.callback(f.fileToLoad)
			}
		} else {
			f.workingDir = f.workingDir + "/../"
			f.scrollVal = 0
		}
	}

	return nil
}
