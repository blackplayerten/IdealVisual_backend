package database

func (db *Database) PqGetKeyFromDetail(detail string) string {
	sms := db.keyParenthesesRe.FindStringSubmatch(detail)
	if len(sms) < 2 {
		return ""
	}

	return sms[1]
}
