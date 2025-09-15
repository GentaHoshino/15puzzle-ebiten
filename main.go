package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text" // textパッケージをインポート
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
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
	emptyX, emptyY int  // 空白の現在位置
	cleared        bool // クリアしたかどうか
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

	// 1...15, 最後が0

	num := 1

	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			if y == g.h-1 && x == g.w-1 {
				g.board[y][x] = 0
			} else {
				g.board[y][x] = num
				num++
			}

		}
	}
	g.emptyX, g.emptyY = g.w-1, g.h-1
}

// JS の処理をGo化：空白(右下)は固定、0..(w*h-2) から2つ選んで swap を n 回
func (g *game) shuffle(n int) {
	max := g.w*g.h - 1 // 右下のインデックスは除外
	for i := 0; i < n; i++ {
		var from, to int
		for from == to {
			from = rand.Intn(max)
			to = rand.Intn(max)
		}

		// 一次元index→（x,y）
		fx, fy := from%g.w, from/g.w
		tx, ty := to%g.w, to/g.w

		g.board[fy][fx], g.board[ty][tx] = g.board[ty][tx], g.board[fy][fx]
	}
}

// 隣接チェック+入れ替え処理
func (g *game) trySwap(x, y int) {
	if abs(g.emptyX-x)+abs(g.emptyY-y) == 1 {
		g.board[g.emptyY][g.emptyX], g.board[y][x] = g.board[y][x], g.board[g.emptyY][g.emptyX]
		g.emptyX, g.emptyY = x, y

	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// クリアしたかどうかを判定するメソッド
func (g *game) isCleared() bool {
	num := 1
	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {

			if y == g.h-1 && x == g.w-1 {
				continue
			}
			if num != g.board[y][x] {
				return false
			}
			num++
		}
	}
	return true
}

func (g *game) Update() error {
	// 既にクリアしていたら、何もしない
	if g.cleared {
		return nil
	}
	// 左クリック
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		gx, gy := mx/TileSize, my/TileSize
		if gx >= 0 && gx < g.w && gy >= 0 && gy < g.h {
			g.trySwap(gx, gy)

			// タイルを動かした後にクリア判定
			if g.isCleared() {
				fmt.Println("クリア！") // コンソールに表示
				g.cleared = true
			}
		}
	}
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

			// もしクリアしていたら色を緑色に変える
			if g.cleared {
				bgColor = color.RGBA{0, 136, 0, 255} // #080
			}
			vector.DrawFilledRect(screen, drawX, drawY, float32(TileSize), float32(TileSize), bgColor, false)

			// タイルの枠線を描画
			borderColor := color.RGBA{255, 136, 0, 255} // #f80
			// ★修正点1: StrokeRectの引数を修正
			vector.StrokeRect(screen, drawX, drawY, float32(TileSize), float32(TileSize), 4, borderColor, false)

			// 数字を描画
			numStr := fmt.Sprintf("%d", val)
			face := basicfont.Face7x13

			// 数字がタイルの中央に来るように位置を計算
			b, _ := font.BoundString(face, numStr)
			w := (b.Max.X - b.Min.X).Ceil() // 幅
			h := (b.Max.Y - b.Min.Y).Ceil() // 高さ

			baseline := (-b.Min.Y).Ceil()
			textX := x*TileSize + (TileSize-w)/2
			textY := y*TileSize + (TileSize-h)/2 + baseline

			text.Draw(screen, numStr, face, textX, textY, color.White)
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano()) // ★ 追加:ランダム化
	ebiten.SetWindowTitle("15パズル")
	// ★修正点2: ウィンドウサイズとnewGameの引数を定数に合わせる
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	g := newGame(BoardWidth, BoardHeight)
	g.shuffle(10000)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
