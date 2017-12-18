package falcona

// Limits
const (
	MaxPly    = 64
	MaxDepth  = 64
	MaxMoves  = 1024
	Checkmate = 0x7FFF - 1 // = 32,766
)

// Game state
const (
	InProgress   = iota
	WhiteWon                         // White checkmated.
	BlackWon                         // Black checkmated.
	Stalemate                        // Draw by stalemate, forced or self-imposed.
	Insufficient                     // Draw by insufficient material.
	Repetition                       // Draw by repetition.
	FiftyMoves                       // Draw by 50 moves rule.
	WhiteWinning = Checkmate / 10    // Decisive advantage for White.
	BlackWinning = -Checkmate/10 + 1 // Decisive advantage for Black.
)

// Piece type
const (
	BlackPawn = iota
	WhitePawn
	BlackKnight
	WhiteKnight
	BlackBishop
	WhiteBishop
	BlackRook
	WhiteRook
	BlackQueen
	WhiteQueen
	BlackKing
	WhiteKing
)

// Color
const (
	White = 0
	Black = 1
)

// Squares
const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// Ranks
const (
	R1 = iota
	R2
	R3
	R4
	R5
	R6
	R7
	R8
)

// Files
const (
	FA = iota
	FB
	FC
	FD
	FE
	FF
	FG
	FH
)

var maskSquare = [64]uint64{
	1 << A1, 1 << B1, 1 << C1, 1 << D1, 1 << E1, 1 << F1, 1 << G1, 1 << H1,
	1 << A2, 1 << B2, 1 << C2, 1 << D2, 1 << E2, 1 << F2, 1 << G2, 1 << H2,
	1 << A3, 1 << B3, 1 << C3, 1 << D3, 1 << E3, 1 << F3, 1 << G3, 1 << H3,
	1 << A4, 1 << B4, 1 << C4, 1 << D4, 1 << E4, 1 << F4, 1 << G4, 1 << H4,
	1 << A5, 1 << B5, 1 << C5, 1 << D5, 1 << E5, 1 << F5, 1 << G5, 1 << H5,
	1 << A6, 1 << B6, 1 << C6, 1 << D6, 1 << E6, 1 << F6, 1 << G6, 1 << H6,
	1 << A7, 1 << B7, 1 << C7, 1 << D7, 1 << E7, 1 << F7, 1 << G7, 1 << H7,
	1 << A8, 1 << B8, 1 << C8, 1 << D8, 1 << E8, 1 << F8, 1 << G8, 1 << H8,
}

var maskRank = [8]uint64{
	0x00000000000000FF, 0x000000000000FF00, 0x0000000000FF0000, 0x00000000FF000000,
	0x000000FF00000000, 0x0000FF0000000000, 0x00FF000000000000, 0xFF00000000000000,
}

var maskFile = [8]uint64{
	0x0101010101010101, 0x0202020202020202, 0x0404040404040404, 0x0808080808080808,
	0x1010101010101010, 0x2020202020202020, 0x4040404040404040, 0x8080808080808080,
}

var maskDiagRight = [15]uint64{
	0x0100000000000000, 0x0201000000000000, 0x0402010000000000, 0x0804020100000000, 0x1008040201000000,
	0x2010080402010000, 0x4020100804020100, 0x8040201008040201, 0x0080402010080402, 0x0000804020100804,
	0x0000008040201008, 0x0000000080402010, 0x0000000000804020, 0x0000000000008040, 0x0000000000000080,
}

var maskDiagLeft = [15]uint64{
	0x0000000000000001, 0x0000000000000102, 0x0000000000010204, 0x0000000001020408, 0x0000000102040810,
	0x0000010204081020, 0x0001020408102040, 0x0102040810204080, 0x0204081020408000, 0x0408102040800000,
	0x0810204080000000, 0x1020408000000000, 0x2040800000000000, 0x4080000000000000, 0x8000000000000000,
}

var maskIsolated = [8]uint64{
	0x0202020202020202, 0x0505050505050505, 0x0A0A0A0A0A0A0A0A, 0x1414141414141414,
	0x2828282828282828, 0x5050505050505050, 0xA0A0A0A0A0A0A0A0, 0x4040404040404040,
}

var castleKingside = [2]uint8{1, 4}
var castleQueenside = [2]uint8{2, 8}
var castleRights = [64]uint8{
	13, 15, 15, 15, 12, 15, 15, 14,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	7, 15, 15, 15, 3, 15, 15, 11,
}

var castleKingsideEmpty = [2]uint64{
	maskSquare[F1] | maskSquare[G1],
	maskSquare[F8] | maskSquare[G8],
}

var castleQueensideEmpty = [2]uint64{
	maskSquare[B1] | maskSquare[C1] | maskSquare[D1],
	maskSquare[B8] | maskSquare[C8] | maskSquare[D8],
}

var castleKingsideSafe = [2]uint64{
	maskSquare[E1] | maskSquare[F1] | maskSquare[G1],
	maskSquare[E8] | maskSquare[F8] | maskSquare[G8],
}

var castleQueensideSafe = [2]uint64{
	maskSquare[C1] | maskSquare[D1] | maskSquare[E1],
	maskSquare[C8] | maskSquare[D8] | maskSquare[E8],
}

var rookMagic = [64]Magic{
	{0x000101010101017E, 0xE580008110204000}, {0x000202020202027C, 0x0160002008011000},
	{0x000404040404047A, 0x3520080020000400}, {0x0008080808080876, 0x6408080448002002},
	{0x001010101010106E, 0x1824080004020400}, {0x002020202020205E, 0x9E04088101020400},
	{0x004040404040403E, 0x8508008001000200}, {0x008080808080807E, 0x5A000040A4840201},
	{0x0001010101017E00, 0x0172800180504000}, {0x0002020202027C00, 0x0094200020483000},
	{0x0004040404047A00, 0x0341400410000800}, {0x0008080808087600, 0x007E040020041200},
	{0x0010101010106E00, 0x0042104404000800}, {0x0020202020205E00, 0x00CC042090008400},
	{0x0040404040403E00, 0x00A88081000A0200}, {0x0080808080807E00, 0x006C810180840100},
	{0x00010101017E0100, 0x0021150010208000}, {0x00020202027C0200, 0x0000158060004000},
	{0x00040404047A0400, 0x0080BC2808402000}, {0x0008080808760800, 0x00113E1000200800},
	{0x00101010106E1000, 0x0002510001010800}, {0x00202020205E2000, 0x0008660814001200},
	{0x00404040403E4000, 0x0000760400408080}, {0x00808080807E8000, 0x0004158180800300},
	{0x000101017E010100, 0x00C0826880004000}, {0x000202027C020200, 0x0000803901002800},
	{0x000404047A040400, 0x0020005D88000400}, {0x0008080876080800, 0x0000136A18001000},
	{0x001010106E101000, 0x0000053200100A00}, {0x002020205E202000, 0x0004301AA0010400},
	{0x004040403E404000, 0x0001085C80048200}, {0x008080807E808000, 0x000000BE04040100},
	{0x0001017E01010100, 0x0001042142401000}, {0x0002027C02020200, 0x0140300036282000},
	{0x0004047A04040400, 0x0009802008801000}, {0x0008087608080800, 0x0000201870841000},
	{0x0010106E10101000, 0x0002020926000A00}, {0x0020205E20202000, 0x0004040078800200},
	{0x0040403E40404000, 0x0002000250900100}, {0x0080807E80808000, 0x00008001B9800100},
	{0x00017E0101010100, 0x0240001090418000}, {0x00027C0202020200, 0x0000100008334000},
	{0x00047A0404040400, 0x0060086000234800}, {0x0008760808080800, 0x0088088011121000},
	{0x00106E1010101000, 0x000204000E1D2800}, {0x00205E2020202000, 0x00020082002A0400},
	{0x00403E4040404000, 0x0006021050270A00}, {0x00807E8080808000, 0x0000610002548080},
	{0x007E010101010100, 0x0401410280201A00}, {0x007C020202020200, 0x0001408100104E00},
	{0x007A040404040400, 0x0020100980084480}, {0x0076080808080800, 0x000A004808286E00},
	{0x006E101010101000, 0x0000100120024600}, {0x005E202020202000, 0x0002000902008E00},
	{0x003E404040404000, 0x0001010800807200}, {0x007E808080808000, 0x0002050400295A00},
	{0x7E01010101010100, 0x004100A080004257}, {0x7C02020202020200, 0x000011008048C005},
	{0x7A04040404040400, 0x8040441100086001}, {0x7608080808080800, 0x8014400C02002066},
	{0x6E10101010101000, 0x00080220010A010E}, {0x5E20202020202000, 0x0002000810041162},
	{0x3E40404040404000, 0x00010486100110E4}, {0x7E80808080808000, 0x000100020180406F},
}

var bishopMagic = [64]Magic{
	{0x0040201008040200, 0x00A08800240C0040}, {0x0000402010080400, 0x0020085008088000},
	{0x0000004020100A00, 0x0080440306000000}, {0x0000000040221400, 0x000C520480000000},
	{0x0000000002442800, 0x0024856000000000}, {0x0000000204085000, 0x00091430A0000000},
	{0x0000020408102000, 0x0009081820C00000}, {0x0002040810204000, 0x0005100400609000},
	{0x0020100804020000, 0x0000004448220400}, {0x0040201008040000, 0x00001030010C0480},
	{0x00004020100A0000, 0x00003228120A0000}, {0x0000004022140000, 0x00004C2082400000},
	{0x0000000244280000, 0x0000108D50000000}, {0x0000020408500000, 0x0000128821100000},
	{0x0002040810200000, 0x00000A0410245000}, {0x0004081020400000, 0x0000030180208800},
	{0x0010080402000200, 0x0050000820500C00}, {0x0020100804000400, 0x0008010050180200},
	{0x004020100A000A00, 0x0150008A00200500}, {0x0000402214001400, 0x000E000505090000},
	{0x0000024428002800, 0x0014001239200000}, {0x0002040850005000, 0x0006000540186000},
	{0x0004081020002000, 0x0001000603502000}, {0x0008102040004000, 0x0005000104480C00},
	{0x0008040200020400, 0x0011400048880200}, {0x0010080400040800, 0x0008080020304800},
	{0x0020100A000A1000, 0x00241800B4000C00}, {0x0040221400142200, 0x0020480010820040},
	{0x0002442800284400, 0x011484002E822000}, {0x0004085000500800, 0x000C050044822000},
	{0x0008102000201000, 0x0010090020401800}, {0x0010204000402000, 0x0004008004A10C00},
	{0x0004020002040800, 0x0044214001005000}, {0x0008040004081000, 0x0008406000188200},
	{0x00100A000A102000, 0x0002043000120400}, {0x0022140014224000, 0x0010280800120A00},
	{0x0044280028440200, 0x0040120A00802080}, {0x0008500050080400, 0x0018282700021000},
	{0x0010200020100800, 0x0008400C00090C00}, {0x0020400040201000, 0x0024200440050400},
	{0x0002000204081000, 0x004420104000A000}, {0x0004000408102000, 0x0020085060000800},
	{0x000A000A10204000, 0x00090428D0011800}, {0x0014001422400000, 0x00000A6248002C00},
	{0x0028002844020000, 0x000140820A000400}, {0x0050005008040200, 0x0140900848400200},
	{0x0020002010080400, 0x000890400C000200}, {0x0040004020100800, 0x0008240492000080},
	{0x0000020408102000, 0x0010104980400000}, {0x0000040810204000, 0x0001480810300000},
	{0x00000A1020400000, 0x0000060410A80000}, {0x0000142240000000, 0x0000000270150000},
	{0x0000284402000000, 0x000000808C540000}, {0x0000500804020000, 0x00002010202C8000},
	{0x0000201008040200, 0x0020A00810210000}, {0x0000402010080400, 0x0080801804828000},
	{0x0002040810204000, 0x0002420010C09000}, {0x0004081020400000, 0x0000080828180400},
	{0x000A102040000000, 0x0000000890285800}, {0x0014224000000000, 0x0000000008234900},
	{0x0028440200000000, 0x00000000A0244500}, {0x0050080402000000, 0x0000022080900300},
	{0x0020100804020000, 0x0000005002301100}, {0x0040201008040200, 0x0080881014040040},
}
