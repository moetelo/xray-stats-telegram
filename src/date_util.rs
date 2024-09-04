pub fn date_or_today(date: String) -> Result<chrono::NaiveDate, chrono::ParseError> {
    match date.as_str() {
        "" => Ok(chrono::Local::now().date_naive()),
        str_date => chrono::NaiveDate::parse_from_str(str_date, "%Y-%m-%d"),
    }
}
