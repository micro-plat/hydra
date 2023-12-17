package errs

//GetDBError 获取DB ERROR 处理行数为0与数据库错误
func GetDBError(row int64, err error) error {
	if err != nil {
		return err
	}
	if row == 0 {
		return New("影响的行数为0:%w", ErrNotExist)
	}
	return nil
}

//GetErrorString 获取错误消息
func GetErrorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
