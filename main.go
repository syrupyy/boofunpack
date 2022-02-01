package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"gopkg.in/ini.v1"
	"howett.net/plist"
)

// SpritePlist format 2 file format
type SpritePlist struct {
	Frames map[string]struct {
		Frame           string `plist:"frame"`
		Offset          string `plist:"offset"`
		Rotated         bool   `plist:"rotated"`
		SourceColorRect string `plist:"sourceColorRect"`
		SourceSize      string `plist:"sourceSize"`
	} `plist:"frames"`
	Metadata struct {
		Format          int    `plist:"format"`
		Size            string `plist:"size"`
		TextureFileName string `plist:"textureFileName"`
	} `plist:"metadata"`
}

// SpritePlist format 3 file format
type SpritePlistFormat3 struct {
	Frames map[string]struct {
		//Aliases          []interface{} `plist:"aliases"`
		Frame           string `plist:"textureRect"`
		Offset          string `plist:"spriteOffset"`
		Rotated         bool   `plist:"textureRotated"`
		SourceColorRect string `plist:"sourceColorRect"`
		SourceSize      string `plist:"sourceSize"`
	} `plist:"frames"`
	Metadata struct {
		Format              int    `plist:"format"`
		//PixelFormat         string `plist:"pixelFormat"`
		//PremultiplyAlpha    bool   `plist:"premultiplyAlpha"`
		//RealTextureFileName string `plist:"realTextureFileName"`
		Size                string `plist:"size"`
		//Smartupdate         string `plist:"smartupdate"`
		TextureFileName     string `plist:"textureFileName"`
	} `plist:"metadata"`
}

// SpritePlistAniinfo file format
type SpritePlistAniinfo struct {
	AnimationList map[string]struct {
		FPS       float64 `plist:"FPS"`
		FrameList []int   `plist:"FrameList"`
	} `plist:"animationlist"`
	FrameList []string `plist:"framelist"`
	Name      string   `plist:"name"`
	Texture   string   `plist:"texture"`
	Type      string   `plist:"type"`
}

var cropSprites bool
var groupByAnimation bool
var closeWhenDone bool
func main() {
	// Load config.ini, create it if it doesn't exist
	cfg, err := ini.Load("config.ini")
    if err != nil {
		ioutil.WriteFile("config.ini", []byte("# Crops sprites to their edges, set to false for original animation size\ncrop_sprites = true\n\n# Splits sprites by animation when possible\ngroup_by_animation = true\n\n# Close the program without prompting the user when done executing\nclose_when_done = false"), 0644)
		cfg, err = ini.Load("config.ini")
		if err != nil {
			exit("Could not make config.ini.")
		}
	}
	cropSprites, _ = cfg.Section("").Key("crop_sprites").Bool()
	groupByAnimation, _ = cfg.Section("").Key("group_by_animation").Bool()
	closeWhenDone, _ = cfg.Section("").Key("close_when_done").Bool()

	// Set file name
	var filename string
	if len(os.Args) == 1 {
		exit("No file specified. Drag a file onto the executable itself or pass a path as a command-line argument.")
	} else {
		filename = os.Args[1]
	}
	if(strings.HasSuffix(filename, ".png")) {
		filename = filename[0:len(filename)-4] + ".plist"
	} else if(strings.HasSuffix(filename, "_aniinfo.plist")) {
		filename = filename[0:len(filename)-4] + ".plist"
	}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		exit("Plist file not found.")
	}
    dir := filepath.Dir(filename) + string(os.PathSeparator)

	// Start ripping the actual plist
	plistFile, err := ioutil.ReadFile(filename)
	if err != nil {
		exit(err.Error())
	}
	var formatCheck map[string]map[string]interface{}
	_, err = plist.Unmarshal(plistFile, &formatCheck)
	if err != nil {
		exit(err.Error())
	}
	var spritePlist SpritePlist
	var separator string
	if formatCheck["metadata"]["format"] == uint64(3) {
		var spritePlistFormat3 SpritePlistFormat3
		_, err = plist.Unmarshal(plistFile, &spritePlistFormat3)
		if err != nil {
			exit(err.Error())
		}
		spritePlist = SpritePlist(spritePlistFormat3)
		separator = ","
	} else {
		_, err = plist.Unmarshal(plistFile, &spritePlist)
		if err != nil {
			exit(err.Error())
		}
		separator = ", "
	}
	src, err := imaging.Open(dir + spritePlist.Metadata.TextureFileName)
	if err != nil {
		exit(err.Error())
	}
	ex, err := os.Executable()
    if err != nil {
		exit(err.Error())
	}
	dirAbsPath := filepath.Dir(ex) + string(os.PathSeparator)
	mainDir := spritePlist.Metadata.TextureFileName[0:len(spritePlist.Metadata.TextureFileName)-len(filepath.Ext(spritePlist.Metadata.TextureFileName))] + string(os.PathSeparator)
	for key, element := range spritePlist.Frames {
		fmt.Println(key)
		rect := strings.Split(strings.ReplaceAll(strings.ReplaceAll(element.Frame, "{", ""), "}", ""), separator)
		var width int
		var height int
		if element.Rotated {
			width, _ = strconv.Atoi(rect[3])
			height, _ = strconv.Atoi(rect[2])
		} else {
			width, _ = strconv.Atoi(rect[2])
			height, _ = strconv.Atoi(rect[3])
		}
		if width < 4 && height < 4 {
			continue
		}
		x, _ := strconv.Atoi(rect[0])
		y, _ := strconv.Atoi(rect[1])
		img := imaging.Crop(src, image.Rect(x, y, x + width, y + height))
		if !strings.Contains(key, "/") {
			key = mainDir + key
		}
		newDir := filepath.Dir(key[0:len(key)-len(filepath.Ext(key))])
		if _, err := os.Stat(dirAbsPath + newDir); os.IsNotExist(err) {
			err = os.MkdirAll(dirAbsPath + newDir, 0644)
			if err != nil {
				exit(err.Error())
			}
		}
		if cropSprites {
			if element.Rotated {
				img = imaging.Rotate90(img)
			}
			err = imaging.Save(img, dirAbsPath + key)
			if err != nil {
				exit(err.Error())
			}
		} else {
			offsetRect := strings.Split(strings.ReplaceAll(strings.ReplaceAll(element.Offset, "{", ""), "}", ""), separator)
			var offsetX int
			var offsetY int
			var flip int
			if element.Rotated {
				offsetX, _ = strconv.Atoi(offsetRect[1])
				offsetY, _ = strconv.Atoi(offsetRect[0])
				flip = 1
			} else {
				offsetX, _ = strconv.Atoi(offsetRect[0])
				offsetY, _ = strconv.Atoi(offsetRect[1])
				flip = -1
			}
			realRect := strings.Split(strings.ReplaceAll(strings.ReplaceAll(element.SourceSize, "{", ""), "}", ""), separator)
			realWidth, _ := strconv.Atoi(realRect[0])
			realHeight, _ := strconv.Atoi(realRect[1])
			dst := imaging.New(realWidth, realHeight, color.NRGBA{0, 0, 0, 0})
			dst = imaging.Paste(dst, img, image.Pt((realWidth - width) / 2 + offsetX, (realHeight - height) / 2 + offsetY * flip))
			if element.Rotated {
				dst = imaging.Rotate90(dst)
			}
			err = imaging.Save(dst, dirAbsPath + key)
			if err != nil {
				exit(err.Error())
			}
		}
	}

	// Group by animation, if configured to do so
	if groupByAnimation {
		plistAniinfoName := strings.Replace(filename, ".plist", "_aniinfo.plist", 1)
		if _, err := os.Stat(plistAniinfoName); os.IsNotExist(err) {
			fmt.Println("Could not find animation list, skipping grouping...")
			exit("Done!")
		}
		plistAniinfoFile, err := ioutil.ReadFile(plistAniinfoName)
		if err != nil {
			exit(err.Error())
		}
		var spritePlistAniinfo SpritePlistAniinfo
		_, err = plist.Unmarshal(plistAniinfoFile, &spritePlistAniinfo)
		if err != nil {
			exit(err.Error())
		}
		var used []string
		for key, element := range spritePlistAniinfo.AnimationList {
			if key == "__all__" || key == "_all" {
				continue
			}
			fmt.Println(key)
			for i, frame := range element.FrameList {
				frameFilename := spritePlistAniinfo.FrameList[frame]
				if strings.Contains(frameFilename, "/") {
					input, err := ioutil.ReadFile(dirAbsPath + frameFilename)
					if err != nil {
						exit(err.Error())
					}
					err = ioutil.WriteFile(dirAbsPath + filepath.Dir(frameFilename) + string(os.PathSeparator) + strings.ReplaceAll(mainDir, string(os.PathSeparator), "_") + strings.ReplaceAll(key, " ", "_") + "_" + fmt.Sprintf("%04d", i + 1) + filepath.Ext(frameFilename), input, 0644)
					if err != nil {
						exit(err.Error())
					}
					used = append(used, frameFilename)
				} else {
					input, err := ioutil.ReadFile(dirAbsPath + mainDir + frameFilename)
					if err != nil {
						exit(err.Error())
					}
					err = ioutil.WriteFile(dirAbsPath + mainDir + strings.ReplaceAll(mainDir, string(os.PathSeparator), "_") + strings.ReplaceAll(key, " ", "_") + "_" + fmt.Sprintf("%04d", i + 1) + filepath.Ext(frameFilename), input, 0644)
					if err != nil {
						exit(err.Error())
					}
					used = append(used, mainDir + frameFilename)
				}
			}
		}
		for _, usedFile := range used {
			err = os.Remove(dirAbsPath + usedFile)
			if err != nil && !os.IsNotExist(err) {
				exit(err.Error())
			}
		}
	}
	exit("Done!")
}

func exit(err string) {
	fmt.Println(err)
	if !closeWhenDone {
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}
	os.Exit(0)
}