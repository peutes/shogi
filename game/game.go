package game

import (
	"shogi/board"
	"shogi/piece"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 600
	BoardMarginX = 200 // 左右のマージン

	BoardMarginY = 20 // 上下のマージン（小さくする）
	CellSize     = 60

	StateNormal   = iota // 通常状態
	StateSelected        // 駒が選択された状態
	StateGameOver        // ゲーム終了状態

	DragNone    = iota // ドラッグなし
	DragBoard          // 盤上の駒をドラッグ中
	DragCapture        // 持ち駒をドラッグ中
)

// 持ち駒エリアの定数
const (
	CaptureAreaWidth  = 120 // 持ち駒エリアを少し広げる
	CaptureAreaMargin = 40  // マージンを広げる
)

// 持ち駒エリアの情報
type CaptureArea struct {
	X, Y, Width, Height int
	Player              piece.Player
}

// ゲームの状態
type GameState struct {
	State          int
	SelectedX      int
	SelectedY      int
	MouseX         int
	MouseY         int
	Dragging       int
	DragPieceType  piece.Type
	DragPieceOwner piece.Player
	Message        string
	ValidMoves     [][2]int // 移動可能なマスの座標リスト
}

// ゲーム管理構造体
type Game struct {
	board         *board.Board
	state         GameState
	font          font.Face
	largeFont     font.Face
	senteCaptures CaptureArea
	goteCaptures  CaptureArea
}

// 新しいゲームを作成
func NewGame(normalFont, largeFont font.Face) *Game {
	boardWidth := board.BoardSize * CellSize
	game := &Game{
		board: board.New(),
		state: GameState{
			State:     StateNormal,
			SelectedX: -1,
			SelectedY: -1,
		},
		font:      normalFont,
		largeFont: largeFont,
		senteCaptures: CaptureArea{
			X:      BoardMarginX + boardWidth + CaptureAreaMargin,
			Y:      BoardMarginY, // 変更
			Width:  CaptureAreaWidth,
			Height: CellSize * 7,
			Player: piece.Sente,
		},
		goteCaptures: CaptureArea{
			X:      BoardMarginX - CaptureAreaWidth - CaptureAreaMargin,
			Y:      BoardMarginY, // 変更
			Width:  CaptureAreaWidth,
			Height: CellSize * 7,
			Player: piece.Gote,
		},
	}
	return game
}

// 座標が持ち駒エリア内かチェック
func (ca *CaptureArea) Contains(x, y int) bool {
	return x >= ca.X && x < ca.X+ca.Width &&
		y >= ca.Y && y < ca.Y+ca.Height
}

// クリックされた持ち駒のインデックスを取得
func (ca *CaptureArea) GetPieceIndex(x, y int) int {
	if !ca.Contains(x, y) {
		return -1
	}
	// クリックされた位置から持ち駒のインデックスを計算
	localY := y - ca.Y
	return localY / CellSize
}

// マウス位置から盤上の座標を計算
func (g *Game) getBoardCoordinates(x, y int) (int, int, bool) {
	boardX := (x - BoardMarginX) / CellSize
	boardY := (y - BoardMarginY) / CellSize

	if boardX >= 0 && boardX < board.BoardSize &&
		boardY >= 0 && boardY < board.BoardSize {
		return boardX, boardY, true
	}
	return -1, -1, false
}

// 持ち駒エリアの座標を取得
func (g *Game) getCaptureCoordinates(x, y int) (int, piece.Player, bool) {
	// 先手の持ち駒エリアをチェック
	if index := g.senteCaptures.GetPieceIndex(x, y); index >= 0 {
		return index, piece.Sente, true
	}

	// 後手の持ち駒エリアをチェック
	if index := g.goteCaptures.GetPieceIndex(x, y); index >= 0 {
		return index, piece.Gote, true
	}

	return -1, piece.None, false
}

// 持ち駒のインデックスから駒の種類を取得
func (g *Game) getPieceTypeFromCaptureIndex(index int, player piece.Player) piece.Type {
	if index < 0 || index >= 7 {
		return piece.Empty
	}

	pieceTypes := []piece.Type{
		piece.Pawn, piece.Lance, piece.Knight,
		piece.Silver, piece.Gold, piece.Bishop, piece.Rook,
	}

	var count int
	if player == piece.Sente {
		count = g.board.SenteCaptures[pieceTypes[index]]
	} else {
		count = g.board.GoteCaptures[pieceTypes[index]]
	}

	if index < len(pieceTypes) && count > 0 {
		return pieceTypes[index]
	}
	return piece.Empty
}

// 持ち駒をドラッグ開始
func (g *Game) startCaptureDrag(x, y int) {
	var targetArea *CaptureArea
	if g.senteCaptures.Contains(x, y) {
		targetArea = &g.senteCaptures
	} else if g.goteCaptures.Contains(x, y) {
		targetArea = &g.goteCaptures
	}

	if targetArea != nil {
		index := targetArea.GetPieceIndex(x, y)
		if index >= 0 {
			captures := g.board.GetCaptures(targetArea.Player)
			if index < len(captures) {
				pieceType := captures[index]
				if targetArea.Player == g.board.CurrentPlayer {
					g.state.Dragging = DragCapture
					g.state.DragPieceType = pieceType
					g.state.DragPieceOwner = targetArea.Player
					// 持ち駒の配置可能な位置を計算
					g.state.ValidMoves = g.board.GetValidDropPositions(pieceType)
				}
			}
		}
	}
}

// 持ち駒の配置を試みる
func (g *Game) tryDropPiece(x, y int) bool {
	boardX, boardY, valid := g.getBoardCoordinates(x, y)
	if !valid {
		return false
	}

	// 移動可能な位置かチェック
	validPosition := false
	for _, pos := range g.state.ValidMoves {
		if pos[0] == boardX && pos[1] == boardY {
			validPosition = true
			break
		}
	}

	if validPosition {
		// 持ち駒を配置
		g.board.DropPiece(boardX, boardY, g.state.DragPieceType)
		return true
	}
	return false
}

// 入力の更新処理
func (g *Game) Update() error {
	// マウス位置の更新
	g.state.MouseX, g.state.MouseY = ebiten.CursorPosition()

	// ゲームオーバー状態の場合
	if g.state.State == StateGameOver {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			// クリックで新しいゲームを開始
			g.board = board.New()
			g.state = GameState{State: StateNormal}
		}
		return nil
	}

	// マウスの入力処理
	g.handleMouseInput()

	return nil
}

// マウス入力の処理
func (g *Game) handleMouseInput() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if g.state.Dragging == DragNone {
			g.handleMousePress()
		}
	} else {
		if g.state.Dragging != DragNone {
			g.handleMouseRelease()
		}
	}
}

// マウスボタン押下時の処理
func (g *Game) handleMousePress() {
	// 盤上の座標を取得
	boardX, boardY, onBoard := g.getBoardCoordinates(g.state.MouseX, g.state.MouseY)

	if onBoard {
		g.handleBoardPress(boardX, boardY)
	} else {
		// 持ち駒エリアの処理
		index, player, inCaptureArea := g.getCaptureCoordinates(g.state.MouseX, g.state.MouseY)
		if inCaptureArea && player == g.board.CurrentPlayer {
			// テスト用に持ち駒の位置を出力
			println("Capture area clicked:", index, player)
			pieceType := g.getPieceTypeFromCaptureIndex(index, player)
			if pieceType != piece.Empty {
				// 持ち駒をドラッグ開始
				g.state.State = StateSelected
				g.state.Dragging = DragCapture
				g.state.DragPieceType = pieceType
				g.state.DragPieceOwner = player
				g.state.ValidMoves = g.board.GetValidDropPositions(pieceType)
				// テスト用に有効な移動位置の数を出力
				println("Valid moves:", len(g.state.ValidMoves))
			}
		}
	}
}

// 盤上でのマウス押下処理
func (g *Game) handleBoardPress(x, y int) {
	p := g.board.GetPiece(x, y)

	if g.state.State == StateNormal {
		// 自分の駒を選択
		if p.Type != piece.Empty && p.Player == g.board.CurrentPlayer {
			g.state.SelectedX = x
			g.state.SelectedY = y
			g.state.State = StateSelected
			g.state.Dragging = DragBoard
			g.state.DragPieceType = p.Type
			g.state.DragPieceOwner = p.Player
			g.calculateValidMoves() // 移動可能なマスを計算
		}
	} else if g.state.State == StateSelected {
		// 移動先の選択
		move := board.Move{
			FromX: g.state.SelectedX,
			FromY: g.state.SelectedY,
			ToX:   x,
			ToY:   y,
		}

		// 同じ位置なら選択解除
		if g.state.SelectedX == x && g.state.SelectedY == y {
			g.resetSelection()
			return
		}

		// 移動処理
		if g.board.IsValidMove(move) {
			g.handleMove(move)
		} else if p.Type != piece.Empty && p.Player == g.board.CurrentPlayer {
			// 無効な移動先なら、新しい駒を選択
			g.state.SelectedX = x
			g.state.SelectedY = y
			g.state.Dragging = DragBoard
			g.state.DragPieceType = p.Type
			g.state.DragPieceOwner = p.Player
		} else {
			g.resetSelection()
		}
	}
}

// マウスボタン解放時の処理
func (g *Game) handleMouseRelease() {
	if g.state.Dragging == DragNone {
		return
	}

	boardX, boardY, onBoard := g.getBoardCoordinates(g.state.MouseX, g.state.MouseY)
	if !onBoard {
		g.resetSelection()
		return
	}

	if g.state.Dragging == DragCapture {
		move := board.Move{
			FromX: -1,
			FromY: -1,
			ToX:   boardX,
			ToY:   boardY,
			Piece: g.state.DragPieceType,
		}

		if g.board.IsValidMove(move) {
			g.board.DropPiece(boardX, boardY, g.state.DragPieceType)
		}
	} else if g.state.Dragging == DragBoard {
		// 盤上の駒の移動処理は変更なし
		move := board.Move{
			FromX: g.state.SelectedX,
			FromY: g.state.SelectedY,
			ToX:   boardX,
			ToY:   boardY,
		}

		if g.board.IsValidMove(move) {
			g.handleMove(move)
		}
	}

	g.resetSelection()
}

// 移動の実行
func (g *Game) handleMove(move board.Move) {
	// 成りの確認
	if g.shouldPromote(move) {
		move.Promote = true
	}

	// 移動を実行
	g.board.MakeMove(move)

	// 王手判定
	if g.board.IsCheck() {
		g.state.Message = "王手！"
		g.state.State = StateGameOver
	} else {
		g.state.Message = ""
	}

	g.resetSelection()
}

// 選択状態のリセット
func (g *Game) resetSelection() {
	g.state.State = StateNormal
	g.state.SelectedX = -1
	g.state.SelectedY = -1
	g.state.Dragging = DragNone
	g.state.ValidMoves = nil // 移動可能なマスをクリア
}

// 成るかどうかの判定
func (g *Game) shouldPromote(move board.Move) bool {
	// 駒打ちの場合は成れない
	if move.FromX == -1 && move.FromY == -1 {
		return false
	}

	p := g.board.GetPiece(move.FromX, move.FromY)
	if !p.Type.CanPromote() {
		return false
	}

	// 強制的な成り
	if g.board.CurrentPlayer == piece.Sente {
		if (p.Type == piece.Pawn || p.Type == piece.Lance) && move.ToY == 0 {
			return true
		}
		if p.Type == piece.Knight && move.ToY <= 1 {
			return true
		}
	} else {
		if (p.Type == piece.Pawn || p.Type == piece.Lance) && move.ToY == board.BoardSize-1 {
			return true
		}
		if p.Type == piece.Knight && move.ToY >= board.BoardSize-2 {
			return true
		}
	}

	// 成れる場所での成り判定
	// 実装では常に成ることにします（実際のゲームではダイアログ等で確認が必要）
	canPromote := false
	if g.board.CurrentPlayer == piece.Sente {
		if move.FromY <= 2 || move.ToY <= 2 {
			canPromote = true
		}
	} else {
		if move.FromY >= 6 || move.ToY >= 6 {
			canPromote = true
		}
	}

	return canPromote
}

// 移動可能なマスを計算
func (g *Game) calculateValidMoves() {
	g.state.ValidMoves = nil
	if g.state.Dragging == DragNone {
		return
	}

	// 盤上の全てのマスをチェック
	for y := 0; y < board.BoardSize; y++ {
		for x := 0; x < board.BoardSize; x++ {
			var move board.Move
			if g.state.Dragging == DragBoard {
				move = board.Move{
					FromX: g.state.SelectedX,
					FromY: g.state.SelectedY,
					ToX:   x,
					ToY:   y,
				}
			} else if g.state.Dragging == DragCapture {
				move = board.Move{
					FromX: -1,
					FromY: -1,
					ToX:   x,
					ToY:   y,
					Piece: g.state.DragPieceType,
				}
			}

			// 移動が有効な場合、座標を追加
			if g.board.IsValidMove(move) {
				if g.state.Dragging == DragBoard &&
					x == g.state.SelectedX && y == g.state.SelectedY {
					continue // 同じ場所は除外
				}
				g.state.ValidMoves = append(g.state.ValidMoves, [2]int{x, y})
			}
		}
	}
}
