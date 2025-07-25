package piece

// 駒の種類
type Type int

const (
	Empty      Type = iota
	Pawn            // 歩
	Lance           // 香
	Knight          // 桂
	Silver          // 銀
	Gold            // 金
	Bishop          // 角
	Rook            // 飛
	King            // 王/玉
	PromPawn        // と金
	PromLance       // 成香
	PromKnight      // 成桂
	PromSilver      // 成銀
	PromBishop      // 馬
	PromRook        // 龍
)

// プレイヤー
type Player int

const (
	None  Player = iota
	Sente        // 先手（下手）
	Gote         // 後手（上手）
)

// 反対のプレイヤーを返す
func (p Player) Opposite() Player {
	switch p {
	case Sente:
		return Gote
	case Gote:
		return Sente
	default:
		return None
	}
}

// 駒の情報
type Piece struct {
	Type   Type
	Player Player
}

// 移動方向を表す構造体
type Direction struct {
	DX, DY   int  // 移動方向
	Repeat   bool // 繰り返し移動可能か
	PromOnly bool // 成駒専用の動きか
}

// 駒の移動可能方向を定義
var movements = map[Type][]Direction{
	Pawn:  {{0, -1, false, false}},
	Lance: {{0, -1, true, false}},
	Knight: { // 桂馬の動き（-2は2マス前を表す）
		{-1, -2, false, false}, // 2マス前、1マス左
		{1, -2, false, false},  // 2マス前、1マス右
	},
	Silver: {
		{-1, -1, false, false},
		{0, -1, false, false},
		{1, -1, false, false},
		{-1, 1, false, false},
		{1, 1, false, false},
	},
	Gold: {
		{-1, -1, false, false},
		{0, -1, false, false},
		{1, -1, false, false},
		{-1, 0, false, false},
		{1, 0, false, false},
		{0, 1, false, false},
	},
	Bishop: {
		{-1, -1, true, false},
		{1, -1, true, false},
		{-1, 1, true, false},
		{1, 1, true, false},
	},
	Rook: {
		{0, -1, true, false},
		{-1, 0, true, false},
		{1, 0, true, false},
		{0, 1, true, false},
	},
	King: {
		{-1, -1, false, false},
		{0, -1, false, false},
		{1, -1, false, false},
		{-1, 0, false, false},
		{1, 0, false, false},
		{-1, 1, false, false},
		{0, 1, false, false},
		{1, 1, false, false},
	},
}

// 成駒の移動方向
var promotedMovements = map[Type][]Direction{
	PromPawn:   movements[Gold],
	PromLance:  movements[Gold],
	PromKnight: movements[Gold], // 成桂は金と同じ動き
	PromSilver: movements[Gold],
	PromBishop: append(
		movements[Bishop],
		Direction{0, -1, false, true},
		Direction{-1, 0, false, true},
		Direction{1, 0, false, true},
		Direction{0, 1, false, true},
	),
	PromRook: append(
		movements[Rook],
		Direction{-1, -1, false, true},
		Direction{1, -1, false, true},
		Direction{-1, 1, false, true},
		Direction{1, 1, false, true},
	),
}

// 成ることができる駒かどうかを判定
func (t Type) CanPromote() bool {
	switch t {
	case Pawn, Lance, Knight, Silver, Bishop, Rook:
		return true
	default:
		return false
	}
}

// 駒の移動可能な方向を取得（後手の場合は方向を反転）
func (p Piece) GetMovements() []Direction {
	// 成り駒の場合は専用の動きを使用
	var dirs []Direction
	if p.Type >= PromPawn && p.Type <= PromRook {
		dirs = promotedMovements[p.Type]
	} else {
		dirs = movements[p.Type]
	}

	if dirs == nil {
		return nil
	}

	if p.Player == Gote {
		// 後手の場合は方向を反転
		reversed := make([]Direction, len(dirs))
		for i, dir := range dirs {
			reversed[i] = Direction{
				DX:       dir.DX,
				DY:       -dir.DY,
				Repeat:   dir.Repeat,
				PromOnly: dir.PromOnly,
			}
		}
		return reversed
	}
	return dirs
}

// 駒の文字列表現を取得（将棋の駒文字）
func (p Piece) String() string {
	switch p.Type {
	case Empty:
		return "　"
	case Pawn:
		return "歩"
	case Lance:
		return "香"
	case Knight:
		return "桂"
	case Silver:
		return "銀"
	case Gold:
		return "金"
	case Bishop:
		return "角"
	case Rook:
		return "飛"
	case King:
		if p.Player == Sente {
			return "玉"
		}
		return "王"
	case PromPawn:
		return "と"
	case PromLance:
		return "杏"
	case PromKnight:
		return "圭"
	case PromSilver:
		return "全"
	case PromBishop:
		return "馬"
	case PromRook:
		return "龍"
	default:
		return "？"
	}
}
