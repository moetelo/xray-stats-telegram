use crate::stats::QueryDate;

pub fn date_or_today(date: String) -> Result<QueryDate, chrono::ParseError> {
    match date.as_str() {
        "" => Ok(QueryDate::now_year_month()),
        str_date => str_date.parse::<QueryDate>(),
    }
}
