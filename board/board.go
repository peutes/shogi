package board

import (
	"shogi/piece"
)

const BoardSize = 9

// 将棋盤の状態を管理する構造体
type Board struct {
	Grid          [BoardSize][BoardSize]piece.Piece
	SenteCaptures map[piece.Type]int
	GoteCaptures  map[piece.Type]int
	CurrentPlayer piece.Player
}

// 移動を表す構造体
type Move struct {
	FromX, FromY int        // 移動元の座標（持ち駒の場合は-1, -1）
	ToX, ToY     int        // 移動先の座標
	Piece        piece.Type // 移動する駒の種類（持ち駒の場合に使用）
	Promote      bool       // 成るかどうか
}

// 新しい将棋盤を初期化
func New() *Board {
	b := &Board{
		SenteCaptures: make(map[piece.Type]int),
		GoteCaptures:  make(map[piece.Type]int),
		CurrentPlayer: piece.Sente,
	}

	// 駒の初期配置
	b.initializePieces()
	return b
}

// 駒の初期配置を設定
func (b *Board) initializePieces() {
	// 先手の駒（下側）
	b.placePiece(8, 0, piece.Lance, piece.Sente)
	b.placePiece(8, 1, piece.Knight, piece.Sente)
	b.placePiece(8, 2, piece.Silver, piece.Sente)
	b.placePiece(8, 3, piece.Gold, piece.Sente)
	b.placePiece(8, 4, piece.King, piece.Sente)
	b.placePiece(8, 5, piece.Gold, piece.Sente)
	b.placePiece(8, 6, piece.Silver, piece.Sente)
	b.placePiece(8, 7, piece.Knight, piece.Sente)
	b.placePiece(8, 8, piece.Lance, piece.Sente)

	b.placePiece(7, 1, piece.Bishop, piece.Sente)
	b.placePiece(7, 7, piece.Rook, piece.Sente)

	for i := 0; i < BoardSize; i++ {
		b.placePiece(6, i, piece.Pawn, piece.Sente)
	}

	// 後手の駒（上側）
	b.placePiece(0, 0, piece.Lance, piece.Gote)
	b.placePiece(0, 1, piece.Knight, piece.Gote)
	b.placePiece(0, 2, piece.Silver, piece.Gote)
	b.placePiece(0, 3, piece.Gold, piece.Gote)
	b.placePiece(0, 4, piece.King, piece.Gote)
	b.placePiece(0, 5, piece.Gold, piece.Gote)
	b.placePiece(0, 6, piece.Silver, piece.Gote)
	b.placePiece(0, 7, piece.Knight, piece.Gote)
	b.placePiece(0, 8, piece.Lance, piece.Gote)

	b.placePiece(1, 1, piece.Rook, piece.Gote)
	b.placePiece(1, 7, piece.Bishop, piece.Gote)

	for i := 0; i < BoardSize; i++ {
		b.placePiece(2, i, piece.Pawn, piece.Gote)
	}
}

// 駒を配置
func (b *Board) placePiece(y, x int, pieceType piece.Type, player piece.Player) {
	b.Grid[y][x] = piece.Piece{Type: pieceType, Player: player}
}

// 指定位置の駒を取得
func (b *Board) GetPiece(x, y int) piece.Piece {
	return b.Grid[y][x]
}

// 移動が有効かチェック
func (b *Board) IsValidMove(move Move) bool {
	// 駒打ちの場合
	if move.FromX == -1 && move.FromY == -1 {
		// 移動先の座標が盤の範囲内かチェック
		if move.ToX < 0 || move.ToX >= BoardSize || move.ToY < 0 || move.ToY >= BoardSize {
			return false
		}
		return b.isValidDrop(move)
	}

	// 通常の移動の場合
	// 移動元と移動先の座標が盤の範囲内かチェック
	if move.FromX < 0 || move.FromX >= BoardSize || move.FromY < 0 || move.FromY >= BoardSize ||
		move.ToX < 0 || move.ToX >= BoardSize || move.ToY < 0 || move.ToY >= BoardSize {
		return false
	}

	return b.isValidNormalMove(move)
}

// 駒打ちの有効性をチェック
func (b *Board) isValidDrop(move Move) bool {
	// 持ち駒があるかチェック
	if b.CurrentPlayer == piece.Sente {
		if b.SenteCaptures[move.Piece] <= 0 {
			return false
		}
	} else {
		if b.GoteCaptures[move.Piece] <= 0 {
			return false
		}
	}

	// 移動先が空いているかチェック
	if b.Grid[move.ToY][move.ToX].Type != piece.Empty {
		return false
	}

	// 歩、香車、桂馬の特殊ルール
	switch move.Piece {
	case piece.Pawn:
		// 二歩のチェック
		if b.hasOwnPawnInColumn(move.ToX) {
			return false
		}
		// 最奥の段には打てない
		if (b.CurrentPlayer == piece.Sente && move.ToY == 0) ||
			(b.CurrentPlayer == piece.Gote && move.ToY == BoardSize-1) {
			return false
		}
	case piece.Lance:
		// 最奥の段には打てない
		if (b.CurrentPlayer == piece.Sente && move.ToY == 0) ||
			(b.CurrentPlayer == piece.Gote && move.ToY == BoardSize-1) {
			return false
		}
	case piece.Knight:
		// 最奥の2段には打てない
		if (b.CurrentPlayer == piece.Sente && move.ToY <= 1) ||
			(b.CurrentPlayer == piece.Gote && move.ToY >= BoardSize-2) {
			return false
		}
	}

	return true
}

// 同じ段に自分の歩があるかチェック
func (b *Board) hasOwnPawnInColumn(x int) bool {
	for y := 0; y < BoardSize; y++ {
		p := b.Grid[y][x]
		if p.Type == piece.Pawn && p.Player == b.CurrentPlayer {
			return true
		}
	}
	return false
}

// 持ち駒を取得
func (b *Board) GetCaptures(player piece.Player) []piece.Type {
	var captures []piece.Type
	var captureMap map[piece.Type]int

	if player == piece.Sente {
		captureMap = b.SenteCaptures
	} else {
		captureMap = b.GoteCaptures
	}

	// 駒の種類の順序を固定
	pieceTypes := []piece.Type{
		piece.Pawn, piece.Lance, piece.Knight,
		piece.Silver, piece.Gold, piece.Bishop, piece.Rook,
	}

	// 固定順序で持ち駒を追加
	for _, pt := range pieceTypes {
		count := captureMap[pt]
		for i := 0; i < count; i++ {
			captures = append(captures, pt)
		}
	}

	return captures
}

// 持ち駒を配置
func (b *Board) DropPiece(x, y int, pieceType piece.Type) {
	// 持ち駒を減らす
	if b.CurrentPlayer == piece.Sente {
		b.SenteCaptures[pieceType]--
	} else {
		b.GoteCaptures[pieceType]--
	}

	// 駒を配置
	b.Grid[y][x] = piece.Piece{Type: pieceType, Player: b.CurrentPlayer}

	// 手番を交代
	b.CurrentPlayer = b.CurrentPlayer.Opposite()
}

// 持ち駒の配置可能な位置を取得
func (b *Board) GetValidDropPositions(pieceType piece.Type) [][2]int {
	var positions [][2]int

	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			move := Move{
				FromX: -1,
				FromY: -1,
				ToX:   x,
				ToY:   y,
				Piece: pieceType,
			}
			if b.isValidDrop(move) {
				positions = append(positions, [2]int{x, y})
			}
		}
	}

	return positions
}

// 通常の移動の有効性をチェック
func (b *Board) isValidNormalMove(move Move) bool {
	p := b.Grid[move.FromY][move.FromX]

	// 移動元の駒が存在し、現在のプレイヤーの駒か確認
	if p.Type == piece.Type(piece.Empty) || p.Player != b.CurrentPlayer {
		return false
	}

	// 移動先に自分の駒がないか確認
	dest := b.Grid[move.ToY][move.ToX]
	if dest.Type != piece.Type(piece.Empty) && dest.Player == b.CurrentPlayer {
		return false
	}

	// 成りの条件をチェック
	if !b.isValidPromotion(move, p) {
		return false
	}

	// 各駒の移動可能範囲をチェック
	return b.isValidPieceMove(move, p)
}

// 成りの有効性をチェック
func (b *Board) isValidPromotion(move Move, p piece.Piece) bool {
	if move.Promote {
		// 成れる駒かチェック
		if !p.Type.CanPromote() {
			return false
		}

		// 敵陣または敵陣から出る場合のみ成れる
		canPromote := false
		if b.CurrentPlayer == piece.Sente {
			if move.FromY <= 2 || move.ToY <= 2 {
				canPromote = true
			}
		} else {
			if move.FromY >= 6 || move.ToY >= 6 {
				canPromote = true
			}
		}

		if !canPromote {
			return false
		}
	}

	// 必ず成らなければならない状況のチェック
	if !move.Promote && b.mustPromote(move, p) {
		return false
	}

	return true
}

// 駒打ちの位置が有効かチェック
func (b *Board) canDropToPosition(move Move) bool {
	// 歩、香車、桂馬の特殊ルール
	if b.CurrentPlayer == piece.Sente {
		switch move.Piece {
		case piece.Pawn, piece.Lance:
			if move.ToY == 0 {
				return false
			}
		case piece.Knight:
			if move.ToY <= 1 {
				return false
			}
		}
	} else {
		switch move.Piece {
		case piece.Pawn, piece.Lance:
			if move.ToY == BoardSize-1 {
				return false
			}
		case piece.Knight:
			if move.ToY >= BoardSize-2 {
				return false
			}
		}
	}
	return true
}

// 必ず成らなければならない状況かチェック
func (b *Board) mustPromote(move Move, p piece.Piece) bool {
	if b.CurrentPlayer == piece.Sente {
		switch p.Type {
		case piece.Pawn, piece.Lance:
			return move.ToY == 0
		case piece.Knight:
			return move.ToY <= 1
		}
	} else {
		switch p.Type {
		case piece.Pawn, piece.Lance:
			return move.ToY == BoardSize-1
		case piece.Knight:
			return move.ToY >= BoardSize-2
		}
	}
	return false
}

// 各駒の移動可能範囲をチェック
func (b *Board) isValidPieceMove(move Move, p piece.Piece) bool {
	dx := move.ToX - move.FromX
	dy := move.ToY - move.FromY

	// 駒の移動可能な方向を取得
	for _, dir := range p.GetMovements() {
		// 桂馬の場合は途中のマスを無視して、移動先が正しいかだけを確認
		if p.Type == piece.Knight {
			if dx == dir.DX && dy == dir.DY {
				return true
			}
			continue
		}

		// その他の駒の場合
		if dir.DX == sign(dx) && dir.DY == sign(dy) {
			if !dir.Repeat {
				// 1マスだけ動く駒の場合
				return abs(dx) == abs(dir.DX) && abs(dy) == abs(dir.DY)
			}
			// 複数マス動ける駒の場合は経路チェック
			if b.isPathClear(move) {
				return true
			}
		}
	}
	return false
}

// 経路上に障害物がないかチェック
func (b *Board) isPathClear(move Move) bool {
	dx := sign(move.ToX - move.FromX)
	dy := sign(move.ToY - move.FromY)
	x, y := move.FromX+dx, move.FromY+dy

	for x != move.ToX || y != move.ToY {
		// 盤の範囲外のチェック
		if x < 0 || x >= BoardSize || y < 0 || y >= BoardSize {
			return false
		}
		if b.Grid[y][x].Type != piece.Empty {
			return false
		}
		x, y = x+dx, y+dy
	}
	return true
}

// 移動を実行
func (b *Board) MakeMove(move Move) {
	// 駒打ちの場合
	if move.FromX == -1 && move.FromY == -1 {
		if b.CurrentPlayer == piece.Sente {
			b.SenteCaptures[move.Piece]--
		} else {
			b.GoteCaptures[move.Piece]--
		}
		b.Grid[move.ToY][move.ToX] = piece.Piece{Type: move.Piece, Player: b.CurrentPlayer}
	} else {
		// 通常の移動
		p := b.Grid[move.FromY][move.FromX]

		// 移動先に相手の駒があれば取る
		if dest := b.Grid[move.ToY][move.ToX]; dest.Type != piece.Empty {
			// 成り駒は元の駒に戻して持ち駒に加える
			capturedType := getOriginalPiece(dest.Type)
			if b.CurrentPlayer == piece.Sente {
				b.SenteCaptures[capturedType]++
			} else {
				b.GoteCaptures[capturedType]++
			}
		}

		// 移動元を空にする
		b.Grid[move.FromY][move.FromX] = piece.Piece{}

		// 移動先に駒を配置（必要に応じて成り）
		if move.Promote {
			p.Type = getPromotedPiece(p.Type)
		}
		b.Grid[move.ToY][move.ToX] = p
	}

	// 手番を交代
	b.CurrentPlayer = getNextPlayer(b.CurrentPlayer)
}

// 王手判定
func (b *Board) IsCheck() bool {
	// 王の位置を探す
	kingX, kingY := b.findKing(b.CurrentPlayer)
	if kingX == -1 {
		return false // 王がない（通常はありえない）
	}

	// 一時的に手番を入れ替えて、相手が王を取れるかチェック
	opponent := getNextPlayer(b.CurrentPlayer)
	tempPlayer := b.CurrentPlayer
	b.CurrentPlayer = opponent

	// 盤上の全ての相手の駒について、王を取れるかチェック
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			p := b.Grid[y][x]
			if p.Type != piece.Empty && p.Player == opponent {
				move := Move{FromX: x, FromY: y, ToX: kingX, ToY: kingY}
				if b.isValidNormalMove(move) {
					b.CurrentPlayer = tempPlayer
					return true
				}
			}
		}
	}

	b.CurrentPlayer = tempPlayer
	return false
}

// 王の位置を探す
func (b *Board) findKing(player piece.Player) (int, int) {
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			p := b.Grid[y][x]
			if p.Type == piece.King && p.Player == player {
				return x, y
			}
		}
	}
	return -1, -1
}

// ユーティリティ関数
func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getNextPlayer(current piece.Player) piece.Player {
	if current == piece.Sente {
		return piece.Gote
	}
	return piece.Sente
}

func getPromotedPiece(t piece.Type) piece.Type {
	switch t {
	case piece.Pawn:
		return piece.PromPawn
	case piece.Lance:
		return piece.PromLance
	case piece.Knight:
		return piece.PromKnight
	case piece.Silver:
		return piece.PromSilver
	case piece.Bishop:
		return piece.PromBishop
	case piece.Rook:
		return piece.PromRook
	default:
		return t
	}
}

func getOriginalPiece(t piece.Type) piece.Type {
	switch t {
	case piece.PromPawn:
		return piece.Pawn
	case piece.PromLance:
		return piece.Lance
	case piece.PromKnight:
		return piece.Knight
	case piece.PromSilver:
		return piece.Silver
	case piece.PromBishop:
		return piece.Bishop
	case piece.PromRook:
		return piece.Rook
	default:
		return t
	}
}
