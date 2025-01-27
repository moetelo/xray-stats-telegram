use chrono::Datelike;

pub enum QueryDate {
    Full(chrono::NaiveDate),
    YearMonth(chrono::NaiveDate),
}

impl QueryDate {
    pub fn now_year_month() -> Self {
        let now = chrono::Utc::now().naive_local().date();
        Self::YearMonth(now)
    }

    fn with_day(&self, day: u32) -> Option<Self> {
        match self {
            Self::Full(date) => Some(Self::Full(*date)),
            Self::YearMonth(date) => date.with_day(day).map(Self::Full),
        }
    }

    pub fn queried_dates(&self) -> Vec<String> {
        match self {
            QueryDate::Full(_) => vec![self.to_string()],

            QueryDate::YearMonth(_) => (1..=31)
                .filter_map(|day| self.with_day(day))
                .map(|qd| qd.to_string())
                .collect(),
        }
    }
}

impl std::str::FromStr for QueryDate {
    type Err = chrono::ParseError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        chrono::NaiveDate::parse_from_str(s, "%Y-%m-%d")
            .map(Self::Full)
            .or_else(|_| chrono::NaiveDate::parse_from_str(s, "%Y-%m").map(Self::YearMonth))
    }
}

impl std::fmt::Display for QueryDate {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            QueryDate::Full(date) => write!(f, "{}", date.format("%Y-%m-%d")),
            QueryDate::YearMonth(date) => write!(f, "{}", date.format("%Y-%m")),
        }
    }
}

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use super::*;

    #[test]
    fn test_query_date_from_str() {
        let date = QueryDate::from_str("2021-01-01").unwrap();
        assert_eq!(date.to_string(), "2021-01-01");

        let date = QueryDate::from_str("2021-01").unwrap();
        assert_eq!(date.to_string(), "2021-01");
    }
}
