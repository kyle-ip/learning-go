package dao

import (
    "database/sql"
    "week02/internal/model"
)

func GetUser(aId int64) (*model.User, error) {
    var userId int64
    var username string
    row := db.QueryRow("select id, username from user where id = ?", aId)
    err := row.Scan(&userId, &username)

    // 如果直接把原始的错误 Warp 后抛出，上层为了判断区分是没有查到数据还是其它错误，会不可避免地与 database/sql 发生耦合。
    if err == sql.ErrNoRows {
    	// 因此把 ErrNoRows 吞掉，替换为一个自定义的错误返回，其它错误 Wrap 返回。
        return nil, ErrNoRows
    }
    return &model.User{
        Id:       userId,
        Username: username,
    }, nil
}
