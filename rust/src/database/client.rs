use crate::models::models::UpdateStatus;
use crate::models::models::Workload;
use futures::StreamExt;
use rusqlite::types::{FromSql, FromSqlError, FromSqlResult, ToSqlOutput, ValueRef};
use rusqlite::{Connection, Error, Result, ToSql};

pub fn create_table_if_not_exist() -> Result<()> {
    let conn = Connection::open("data.db")?;
    conn.execute(
        "CREATE TABLE IF NOT EXISTS workloads (
                  id              INTEGER PRIMARY KEY,
                  name            TEXT NOT NULL,
                  image           TEXT NOT NULL,
                  namespace       TEXT NOT NULL,
                  git_ops_repo    TEXT,
                  include_pattern TEXT,
                  exclude_pattern TEXT,
                  update_available TEXT,
                  current_version TEXT NOT NULL,
                  latest_version  TEXT NOT NULL,
                  last_scanned    TEXT NOT NULL,
                  scan_id        INTEGER,
                  scan_type     TEXT
                  )",
        [],
    )?;
    Ok(())
}

pub fn return_workload(name: String, namespace: String) -> (Result<Workload>) {
    let conn = Connection::open("data.db")?;
    let mut stmt = conn.prepare("SELECT * FROM workloads WHERE name = ?1 AND namespace = ?2")?;
    let mut workload = stmt.query_map(&[&name, &namespace], |row| {
        Ok(Workload {
            name: row.get(1)?,
            image: row.get(2)?,
            namespace: row.get(3)?,
            git_ops_repo: row.get(4)?,
            include_pattern: row.get(5)?,
            exclude_pattern: row.get(6)?,
            update_available: row.get(7)?,
            current_version: row.get(8)?,
            latest_version: row.get(9)?,
            last_scanned: row.get(10)?,
        })
    })?;
    if let Some(workload) = workload.next() {
        return Ok(workload?);
    } else {
        return Err(rusqlite::Error::QueryReturnedNoRows);
    };
}

pub fn return_all_workloads() -> (Result<Vec<Workload>>) {
    let conn = Connection::open("data.db")?;
    // only get latest unique name and namespace combinations
    let mut stmt = conn.prepare("SELECT * FROM workloads WHERE scan_id IN (SELECT MAX(scan_id) FROM workloads GROUP BY scan_type)")?;
    let workloads = stmt.query_map([], |row| {
        Ok(Workload {
            name: row.get(1)?,
            image: row.get(2)?,
            namespace: row.get(3)?,
            git_ops_repo: row.get(4)?,
            include_pattern: row.get(5)?,
            exclude_pattern: row.get(6)?,
            update_available: row.get(7)?,
            current_version: row.get(8)?,
            latest_version: row.get(9)?,
            last_scanned: row.get(10)?,
        })
    })?;
    let mut result = Vec::new();
    for workload in workloads {
        result.push(workload?);
    }
    Ok(result)
}

impl ToSql for UpdateStatus {
    fn to_sql(&self) -> rusqlite::Result<ToSqlOutput<'_>> {
        Ok(self.to_string().into())
    }
}

impl FromSql for UpdateStatus {
    fn column_result(value: ValueRef<'_>) -> FromSqlResult<Self> {
        value
            .as_str()?
            .parse()
            .map_err(|e| FromSqlError::Other(Box::new(e)))
    }
}

pub fn get_latest_scan_id() -> std::result::Result<i32, Error> {
    let conn = Connection::open("data.db")?;
    let mut stmt = conn.prepare("SELECT MAX(scan_id) FROM workloads")?;
    let scan_id_iter = stmt.query_map([], |row| row.get(0))?;

    for scan_id_result in scan_id_iter {
        return scan_id_result.map(|id: Option<i32>| id.unwrap_or(0)); // Handle potential NULL
    }
    Ok(0)
}

pub fn insert_workload(workload: &Workload, scan_id: i32) -> Result<()> {
    let conn = Connection::open("data.db")?;
    //get scan_id
    match conn.execute(
        "INSERT INTO workloads (name, image, namespace, git_ops_repo, include_pattern, exclude_pattern, update_available, current_version, latest_version, last_scanned, scan_id, scan_type)
                  VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12)",
        &[
            &workload.name,
            &workload.image,
            &workload.namespace,
            workload.git_ops_repo.as_ref().map(String::as_str).unwrap_or_default(), // Handle potential None
            workload.include_pattern.as_ref().map(String::as_str).unwrap_or_default(),
            workload.exclude_pattern.as_ref().map(String::as_str).unwrap_or_default(),
            &workload.update_available.to_string(), // Consider an enum for clarity
            &workload.current_version,
            &workload.latest_version,
            &workload.last_scanned,
            &scan_id.to_string(),
            "placeholder",
        ],
    ) {
        Ok(_) => Ok(()),
        Err(e) => Err(e),
    }
}