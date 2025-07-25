package main

import (
	"log"
	"shogi/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// mainパッケージは将棋ゲームのエントリーポイントです。
// Ebitenを使用してゲームを実行します。

// Ebitenのゲームループを開始するために必要な関数を定義します。
var (
	// ゲームの幅と高さ
	gameWidth  = 900
	gameHeight = 600

	// セルのサイズ
	CellSize = 60

	// 将棋盤のマージン
	BoardMarginX = 50
	BoardMarginY = 50
)

// フォントの初期化
func initFont() (font.Face, font.Face) {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	// 通常サイズのフォント
	normalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 大きいサイズのフォント
	largeFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return normalFont, largeFont
}

func main() {
	// ウィンドウ設定
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("将棋")

	// フォントの初期化とゲームの作成
	normalFont, largeFont := initFont()
	g := game.NewGame(normalFont, largeFont)

	// ゲーム開始
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
