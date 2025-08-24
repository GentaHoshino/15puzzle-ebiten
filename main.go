package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text" // textパッケージをインポート
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	BoardWidth   = 4
	BoardHeight  = BoardWidth
	TileSize     = 80
	ScreenWidth  = BoardWidth * TileSize
	ScreenHeight = BoardHeight * TileSize
)

type game struct {
	board          [][]int // 2次元スライス： board[y][x]
	w, h           int
	emptyX, emptyY int // 空白の現在位置
}

func newGame(w, h int) *game {
	g := &game{w: w, h: h}
	g.initBoard()
	return g
}

func (g *game) initBoard() {
	g.board = make([][]int, g.h)
	for y := 0; y < g.h; y++ {
		g.board[y] = make([]int, g.w)
	}

	nums := make([]int, g.w*g.h)
	for i := 0; i < g.w*g.h-1; i++ {
		nums[i] = i + 1
	}
	nums[g.w*g.h-1] = 0

	idx := 0
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			g.board[y][x] = nums[idx]
			idx++
		}
	}
	g.emptyX, g.emptyY = g.w-1, g.h-1
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{50, 50, 50, 255})

	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			drawX, drawY := float32(x*TileSize), float32(y*TileSize)
			val := g.board[y][x]
			if val == 0 {
				// 空白も薄い背景＋細枠（任意）
				vector.DrawFilledRect(screen, drawX, drawY, float32(TileSize), float32(TileSize),
					color.RGBA{60, 60, 60, 255}, true)
				vector.StrokeRect(screen, drawX, drawY, float32(TileSize), float32(TileSize),
					2, color.RGBA{90, 90, 90, 255}, true)
				continue
			}

			// タイルの背景を描画
			bgColor := color.RGBA{136, 68, 0, 255} // #840
			vector.DrawFilledRect(screen, drawX, drawY, float32(TileSize), float32(TileSize), bgColor, false)

			// タイルの枠線を描画
			borderColor := color.RGBA{255, 136, 0, 255} // #f80
			// ★修正点1: StrokeRectの引数を修正
			vector.StrokeRect(screen, drawX, drawY, float32(TileSize), float32(TileSize), 4, borderColor, false)

			// 数字を描画
			numStr := fmt.Sprintf("%d", val)
			fontFace := basicfont.Face7x13

			// 数字がタイルの中央に来るように位置を計算
			bounds := text.BoundString(fontFace, numStr)
			w := bounds.Dx() // 幅（ピクセル）
			h := bounds.Dy() // 高さ（ピクセル）
			textX := x*TileSize + (TileSize-w)/2
			textY := y*TileSize + (TileSize-h)/2 + h

			text.Draw(screen, numStr, fontFace, textX, textY, color.White)
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowTitle("15パズル")
	// ★修正点2: ウィンドウサイズとnewGameの引数を定数に合わせる
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	g := newGame(BoardWidth, BoardHeight)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
