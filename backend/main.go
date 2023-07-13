package main

import(
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
)

type WebRes struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}

// 先把所有棋局和玩家存在内存中
var AllChessGame []*ChessGame
var AllUsers []User
// 棋局编号，会递增
var ChessGameId int

// 棋子
type ChessPiece struct{
	row int
	col int
	val int
}
func NewUser(id int, name string) (user User){
	user = User{
		Id : id,
		Name : name,
	}
	return
}
type User struct {
	Id int
	Name string
}
// 获取玩家信息
func GetUserInfo(userId int) (user User, err error) {
	// todo:从数据库中读取玩家信息
	return user, nil
}

// 获取棋局
func GetChessGame(userId int, chessGameId int) (*ChessGame, error) {
	// var chessGame ChessGame
	// todo:从数据库中读取该棋局
	for i:=0; i<len(AllChessGame); i++ {
		if AllChessGame[i].Id == chessGameId {
			return AllChessGame[i], nil
		}
	}
	err := fmt.Errorf("该棋局不存在")
	return nil, err
}

func NewChessGame(userId int, radius int, rowMax int, colMax int) (chessGame *ChessGame){
	// todo:存储数据库中，获得实际的id
	// 假数据
	ChessGameId++
	id := ChessGameId

	chessGame = &ChessGame{
		Id:id,
		ChessBoardRadius:radius,
		ChessBoardRowMax:rowMax,
		ChessBoardColMax:colMax,
		TotalPiece:rowMax*colMax,
	}
	// 定义先手
	chessGame.Players[0] = userId
	chessGame.Map2Array = make([][]int, rowMax)
	for i:=0; i<len(chessGame.Map2Array); i++ {
		chessGame.Map2Array[i] = make([]int, colMax)
	}
	return chessGame
}
// 棋局
type ChessGame struct {
	// 棋局唯一id
	Id int
	// 半径（五子棋为4）
	ChessBoardRadius int
	// 棋盘边界横向最大值
	ChessBoardRowMax int
	// 棋盘边界纵向最大值
	ChessBoardColMax int
	// 棋盘落子情况-二维数组
	Map2Array [][]int
	// 棋盘落子情况-稀疏数组
	MapSparseArray []ChessPiece
	// 棋格总数
	TotalPiece int
	// 当前对弈操作次数（偶数对应到先手下）
	CurOperations int
	// 玩家id
	Players [2]int
	// 赢家id
	WinnerId int
	IsSecondHand bool
	// 该局棋耗时统计
	ToTalTime int
}

/**
初始化该棋盘落子情况-二维数组
*/
func (this *ChessGame) MapSparseArray2Map2Array() {
	for i:=0; i<len(this.MapSparseArray); i++ {
		x := this.MapSparseArray[i].row
		y := this.MapSparseArray[i].col
		v := this.MapSparseArray[i].val
		this.Map2Array[x][y] = v
	}
}
func (this *ChessGame) Map2Array2MapSparseArray() {

}

/**
获得棋盘落子情况-稀疏数组
*/
func (this *ChessGame) GetChessMapSparseArray() []ChessPiece{
	return this.MapSparseArray
}

/**
判断当前是否先手下
*/
func (this *ChessGame) IsFirstHandTime() bool{
	if this.CurOperations%2 == 0 {
		// 偶数，先手下
		return true
	}
	return false
}
/**
获得当前玩家的棋子值
*/
func (this *ChessGame) GetChessPieceVal(userId int) (val int){
	val = userId
	// if this.IsFirstHandTime() {
	// 	val = 1
	// } else {
	// 	val = 2
	// }
	return
}
/**
判断是否该当前玩家下棋
*/
func (this *ChessGame) IsRightPlayers(playerId int) bool {
	if this.IsFirstHandTime() {
		if playerId == this.Players[0] {
			return true
		} 
	} else {
		if playerId == this.Players[1] {
			return true
		}
	}
	return false
}

/**
添加对弈操作次数
*/
func (this *ChessGame) AddOperations(){
	this.CurOperations++
}

/**
落一步棋
*/
func (this *ChessGame) TakeAMove(userId int, x int, y int) (bool){
	if (this.Map2Array[x][y] != 0) {
		// 该棋格中已有棋子
		return false
	}
	val := this.GetChessPieceVal(userId)
	this.Map2Array[x][y] = val
	this.AddOperations()
	return true
}

/**
落一步棋
*/
func (this *ChessGame) TakeAMoveOld(onePiece ChessPiece) (res bool){
	if (this.Map2Array[onePiece.row][onePiece.col] == 0) {
		// 该棋格中已有棋子
		return
	}
	this.Map2Array[onePiece.row][onePiece.col] = onePiece.val
	this.AddOperations()
	res = true
	return
}
/**
获取能判断赢的范围
*/
func (this *ChessGame) getWinRange(x, y int) (xMin, xMax, yMin, yMax int) {
	// xMin = x - this.ChessBoardRadius
	// if (xMin < 0) {
	// 	xMin = 0
	// }
	if x < this.ChessBoardRadius {
		xMin = 0
	} else {
		xMin = x - this.ChessBoardRadius
	}
	xMax = x + this.ChessBoardRadius
	if (xMax > this.ChessBoardRowMax) {
		xMax = this.ChessBoardRowMax
	}
	// yMin = y - this.ChessBoardRadius
	// if (yMin < 0) {
	// 	yMin = 0
	// }
	if y < this.ChessBoardRadius {
		yMin = 0
	} else {
		yMin = y - this.ChessBoardRadius
	}
	yMax = y + this.ChessBoardRadius
	if (yMax > this.ChessBoardColMax) {
		yMax = this.ChessBoardColMax
	}
	return
}
/**
获取下棋结果
@return: 1赢了；2平局；0未赢；
*/
func (this *ChessGame) GetResult(userId, x, y int) (res int){
	val := this.GetChessPieceVal(userId)
	xMin, xMax, yMin, yMax := this.getWinRange(x, y)
	fmt.Printf("xMin, xMax, yMin, yMax val : %d %d %d %d %d \n", xMin, xMax, yMin, yMax, val)
	num := 0
	// -判断
	for i:=xMin; i<=xMax; i++ {
		if (this.Map2Array[i][y] == val) {
			num++
		}else {
			num = 0
		}
		if (num == 5) {
			res = 1
			return
		}
	}
	num = 0
	// |判断
	for j:=yMax; j>=yMin; j-- {
		if (this.Map2Array[x][j] == val) {
			num++
		}else {
			num = 0
		}
		fmt.Println(num)
		if (num == 5) {
			res = 1
			return
		}
	}
	num = 0
	// \判断
	for i,j:=xMin,yMax; i<=xMax && j>=yMin; i,j=i+1,j-1 {
		if (this.Map2Array[i][j] == val) {
			num++
		}else {
			num = 0
		}
		if (num == 5) {
			res = 1
			return
		}
	}
	num = 0
	// /判断
	for i,j:=xMin,yMin; i<=xMax || j<=yMax; i,j=i+1,j+1 {
		if (this.Map2Array[i][j] == val) {
			num++
		}else {
			num = 0
		}
		if (num == 5) {
			res = 1
			return
		}
	}
	// 判断是否平局
	if this.CurOperations == this.TotalPiece {
		res = 2
	}
	return
}

// 创建棋局
func DoCreatChessGame(userId int) (*ChessGame, error) {
	var chessGame *ChessGame
	chessGame = NewChessGame(userId, 4, 20, 20)
	// 存储
	AllChessGame = append(AllChessGame, chessGame)
	return chessGame, nil
}
// 加入棋局
func DoJoinChessGame(userId int, chessGameId int) (*ChessGame, error) {
	var chessGame *ChessGame
	// todo:从数据库中读取该棋局
	for i:=0; i<len(AllChessGame); i++ {
		if AllChessGame[i].Id == chessGameId {
			AllChessGame[i].Players[1] = userId
			return AllChessGame[i], nil
		}
	}
	err := fmt.Errorf("该棋局不存在")
	return chessGame, err
}
// 下棋
func DoChess(userId int, chessGameId int, x int, y int) (int, error){
	var res int
	// // 判断玩家是否存在
	// _, err := GetUserInfo(userId)
	// if err!=nil {
	// 	err = fmt.Errorf("玩家%d不存在", userId)
	// 	return res, err
	// }
	
	// 获取棋局（chessGameId=0 时新建棋局）
	chessGame, err := GetChessGame(userId, chessGameId)
	if err!=nil {
		return res, err
	}

	// 判断是否到该玩家落子
	if !chessGame.IsRightPlayers(userId) {
		err = fmt.Errorf("当前不到%d落子", userId)
		return res, err
	}
	// 在棋盘上落子
	fmt.Println("落子前：")
	for _, v1 := range chessGame.Map2Array {
		for _, v2 := range v1 {
			fmt.Printf("%d \t", v2)
		}
		fmt.Println()
	}
	if !chessGame.TakeAMove(userId, x, y) {
		err = fmt.Errorf("落子失败，请重试")
		return res, err
	}
	fmt.Println("打印二维数组地图如下：")
	for _, v1 := range chessGame.Map2Array {
		for _, v2 := range v1 {
			fmt.Printf("%d \t", v2)
		}
		fmt.Println()
	}
	// 判断下棋结果，返回
	return chessGame.GetResult(userId, x, y), nil
}

func main() {
	fmt.Println("init ok")
	// 注册路由处理函数
	http.HandleFunc("/chessgame/create", ApiCreatChessGame)
	http.HandleFunc("/chessgame/join", ApiJoinChessGame)
	http.HandleFunc("/chessgame/do", ApiDoChess)

	// 启动Web服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// 处理路由的函数
func ApiCreatChessGame(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}
	chessGame, err := DoCreatChessGame(userId)
	msg := fmt.Sprintf("创建成功，棋局号是 %d ", chessGame.Id)
	webRes := WebRes{
		Code : 0,
		Msg : msg,
	}
	if err != nil {
		webRes = WebRes{
			Code : -1,
			Msg : "fail",
		}
	}
	// 将 person 转换为 JSON 格式
	jsonData, err := json.Marshal(webRes)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 设置响应头为 JSON 类型
	w.Header().Set("Content-Type", "application/json")

	// 将 JSON 数据写入响应体
	w.Write(jsonData)
}
func ApiJoinChessGame(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	chessGameIdStr := r.URL.Query().Get("chessGameId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}
	chessGameId, err := strconv.Atoi(chessGameIdStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}

	_, err = DoJoinChessGame(userId, chessGameId)
	webRes := WebRes{
		Code : 0,
		Msg : "恭喜你，成功加入棋局号 1",
	}
	if err != nil {
		webRes = WebRes{
			Code : -1,
			Msg : "fail",
		}
	}
	// 将 person 转换为 JSON 格式
	jsonData, err := json.Marshal(webRes)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 设置响应头为 JSON 类型
	w.Header().Set("Content-Type", "application/json")

	// 将 JSON 数据写入响应体
	w.Write(jsonData)
}

func ApiDoChess(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("userId")
	chessGameIdStr := r.URL.Query().Get("chessGameId")
	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}
	chessGameId, err := strconv.Atoi(chessGameIdStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}
	x, err := strconv.Atoi(xStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return
	}

	res, err := DoChess(userId, chessGameId, x, y)
	var result string
	if res == 1 {
		result = "恭喜你，你赢了！"
	} else if res == 2 {
		result = "是平局哟！"
	} else {
		result = "还未分出胜负，请继续下棋。"
	}
 	webRes := WebRes{
		Code : 0,
		Msg : result,
	}
	if err != nil {
		fmt.Println("err=", err)
		webRes = WebRes{
			Code : -1,
			Msg : "fail",
		}
	}
	// 将 person 转换为 JSON 格式
	jsonData, err := json.Marshal(webRes)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 设置响应头为 JSON 类型
	w.Header().Set("Content-Type", "application/json")

	// 将 JSON 数据写入响应体
	w.Write(jsonData)
}