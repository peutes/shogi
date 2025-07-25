package game

import (
	"fmt"
	"image/color"
	"shogi/board"
	"shogi/piece"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// 画面描画
func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を描画
	screen.Fill(color.RGBA{220, 220, 220, 255})

	// 将棋盤を描画
	g.drawBoard(screen)

	// 駒を描画
	g.drawPieces(screen)

	// 持ち駒エリアを描画
	g.drawCaptureAreas(screen)

	// UI要素を描画
	g.drawUI(screen)
}

// 将棋盤を描画
func (g *Game) drawBoard(screen *ebiten.Image) {
	boardWidth := board.BoardSize * CellSize
	boardHeight := board.BoardSize * CellSize

	// 盤の背景
	ebitenutil.DrawRect(screen,
		float64(BoardMarginX),
		float64(BoardMarginY),
		float64(boardWidth),
		float64(boardHeight),
		color.RGBA{210, 180, 140, 255})

	// 移動可能なマスをハイライト表示
	for _, pos := range g.state.ValidMoves {
		ebitenutil.DrawRect(screen,
			float64(BoardMarginX+pos[0]*CellSize),
			float64(BoardMarginY+pos[1]*CellSize),
			float64(CellSize),
			float64(CellSize),
			color.RGBA{0, 255, 0, 64})
	}

	// マス目を描画
	for i := 0; i <= board.BoardSize; i++ {
		// 横線
		ebitenutil.DrawLine(screen,
			float64(BoardMarginX),
			float64(BoardMarginY+i*CellSize),
			float64(BoardMarginX+boardWidth),
			float64(BoardMarginY+i*CellSize),
			color.Black)

		// 縦線
		ebitenutil.DrawLine(screen,
			float64(BoardMarginX+i*CellSize),
			float64(BoardMarginY),
			float64(BoardMarginX+i*CellSize),
			float64(BoardMarginY+boardHeight),
			color.Black)
	}
}

// 盤上の駒を描画
func (g *Game) drawPieces(screen *ebiten.Image) {
	for y := 0; y < board.BoardSize; y++ {
		for x := 0; x < board.BoardSize; x++ {
			// 選択された駒のハイライト
			if g.state.State == StateSelected &&
				x == g.state.SelectedX &&
				y == g.state.SelectedY &&
				g.state.Dragging == DragNone {
				ebitenutil.DrawRect(screen,
					float64(BoardMarginX+x*CellSize),
					float64(BoardMarginY+y*CellSize),
					float64(CellSize),
					float64(CellSize),
					color.RGBA{255, 255, 0, 128})
			}

			// 駒の描画（ドラッグ中の駒は除く）
			p := g.board.GetPiece(x, y)
			if p.Type != piece.Empty &&
				!(g.state.Dragging == DragBoard &&
					x == g.state.SelectedX &&
					y == g.state.SelectedY) {
				g.drawPiece(screen, p,
					BoardMarginX+x*CellSize+CellSize/2,
					BoardMarginY+y*CellSize+CellSize/2)
			}
		}
	}
}

// 持ち駒エリアを描画
func (g *Game) drawCaptureAreas(screen *ebiten.Image) {
	// 先手の持ち駒エリア
	g.drawCaptureArea(screen, &g.senteCaptures)
	// 後手の持ち駒エリア
	g.drawCaptureArea(screen, &g.goteCaptures)

	// ドラッグ中の駒を描画
	if g.state.Dragging != DragNone {
		p := piece.Piece{
			Type:   g.state.DragPieceType,
			Player: g.state.DragPieceOwner,
		}
		g.drawPiece(screen, p, g.state.MouseX, g.state.MouseY)
	}
}

// 特定の持ち駒エリアを描画
func (g *Game) drawCaptureArea(screen *ebiten.Image, area *CaptureArea) {
	// エリアの背景
	ebitenutil.DrawRect(screen,
		float64(area.X),
		float64(area.Y),
		float64(area.Width),
		float64(area.Height),
		color.RGBA{230, 220, 210, 255})

	// "持駒" のラベルを描画
	text.Draw(screen,
		"持駒",
		g.font,
		area.X+area.Width/2-20,
		area.Y+25,
		color.Black)

	// 持ち駒の描画
	captures := g.board.GetCaptures(area.Player)
	for i, pieceType := range captures {
		// ドラッグ中の持ち駒は表示しない
		if g.state.Dragging == DragCapture &&
			g.state.DragPieceType == pieceType &&
			g.state.DragPieceOwner == area.Player {
			continue
		}

		p := piece.Piece{
			Type:   pieceType,
			Player: area.Player,
		}
		g.drawPiece(screen, p,
			area.X+area.Width/2,
			area.Y+40+i*CellSize)
	}
}

// UI要素を描画
func (g *Game) drawUI(screen *ebiten.Image) {
	// 手番表示の背景を描画
	ebitenutil.DrawRect(screen,
		float64(ScreenWidth/2-60),
		float64(BoardMarginY-5),
		120,
		25,
		color.RGBA{230, 230, 230, 255})

	// 手番表示
	playerText := "先手番"
	if g.board.CurrentPlayer == piece.Gote {
		playerText = "後手番"
	}
	bounds := text.BoundString(g.font, playerText)
	text.Draw(screen, playerText, g.font,
		ScreenWidth/2-bounds.Dx()/2, // 中央揃え
		BoardMarginY+15,             // 将棋盤の上部に表示
		color.Black)

	// メッセージ表示
	if g.state.Message != "" {
		text.Draw(screen, g.state.Message, g.largeFont,
			ScreenWidth/2-50, 570,
			color.RGBA{255, 0, 0, 255})
	}

	// ドラッグ中の駒を描画
	if g.state.Dragging != DragNone {
		p := piece.Piece{
			Type:   g.state.DragPieceType,
			Player: g.state.DragPieceOwner,
		}
		g.drawPiece(screen, p, g.state.MouseX, g.state.MouseY)
	}

	// ゲームオーバー表示
	if g.state.State == StateGameOver {
		ebitenutil.DrawRect(screen,
			200, 200, 400, 200,
			color.RGBA{0, 0, 0, 200})
		text.Draw(screen, "ゲーム終了", g.largeFont,
			320, 250,
			color.White)
		text.Draw(screen, "クリックで再開", g.font,
			320, 300,
			color.White)
	}
}

// 個々の駒を描画
func (g *Game) drawPiece(screen *ebiten.Image, p piece.Piece, centerX, centerY int) {
	if p.Type == piece.Empty {
		return
	}

	// 駒の背景（丸）
	ebitenutil.DrawCircle(screen,
		float64(centerX),
		float64(centerY),
		float64(CellSize/2-5),
		color.RGBA{240, 215, 160, 255})

	// 文字色の設定（黒または赤）
	textColor := color.RGBA{0, 0, 0, 255} // 黒
	if p.Player == piece.Gote {
		textColor = color.RGBA{200, 0, 0, 255} // 赤
	}

	// 駒の文字を描画（中央揃え）
	pieceStr := p.String()
	bounds := text.BoundString(g.font, pieceStr)
	text.Draw(screen, pieceStr, g.font,
		centerX-bounds.Dx()/2,
		centerY+bounds.Dy()/2,
		textColor)
}

// 持ち駒を描画
func (g *Game) drawCapturedPieces(screen *ebiten.Image, player piece.Player,
	pieceTypes []piece.Type, area CaptureArea) {

	for i, pt := range pieceTypes {
		var count int
		if player == piece.Sente {
			count = g.board.SenteCaptures[pt]
		} else {
			count = g.board.GoteCaptures[pt]
		}

		p := piece.Piece{Type: pt, Player: player}
		// Y位置を調整して、上部の"持駒"テキストの下から開始
		yPos := area.Y + (i+1)*CellSize + 30

		if count > 0 {
			g.drawPiece(screen, p,
				area.X+35, // X位置を左に調整
				yPos)
			text.Draw(screen,
				fmt.Sprintf("x%d", count),
				g.font,
				area.X+65, // カウント表示のX位置も調整
				yPos,
				color.Black)
		} else {
			// 持ち駒がない場合は薄く表示
			text.Draw(screen,
				p.String(),
				g.font,
				area.X+35, // X位置を左に調整
				yPos,
				color.RGBA{150, 150, 150, 255})
		}
	}
}

// レイアウト設定
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
