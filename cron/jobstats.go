package cron

import (
	"database/sql"
	"fmt"
)

type JobStats struct {
	Category	string
	Status		string
	Count		int
}

//jobs/submitted
func GetSubmittedJobCounts(db *sql.DB, startDate string, endDate string)([]JobStats, error) {
	//skip curl wrapper jobs
	query := `SELECT b.job_type, count(b.*) AS count
            FROM (
              SELECT a.id,
                CASE WHEN array_length(a.types, 1) = 1 THEN a.types[1]
              ELSE 'Mixed'
           END AS job_type
           FROM (
               SELECT j.id, array_agg(DISTINCT t.name) AS types FROM jobs j
               JOIN job_steps s ON j.id = s.job_id
               JOIN job_types t ON s.job_type_id = t.id
               WHERE j.start_date >= ($1 :: DATE) AND j.start_date <= ($2 :: DATE) + INTERVAL '1 day'
               AND j.app_id != '1e8f719b-0452-4d39-a2f3-8714793ee3e6'
               GROUP BY j.id
           ) AS a
          ) AS b
         GROUP BY b.job_type
         ORDER BY b.job_type;`

	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []JobStats

	for i := 0; rows.Next(); i++ {
		var job JobStats
		err := rows.Scan(&job.Category, &job.Status, &job.Count)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
		output := fmt.Sprintf("Total no.of %[1]v jobs Submitted in last 24 hours: %[2]v", job.Category, job.Count)
		fmt.Println(output)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

//jobs/status
func GetJobStatusCounts(db *sql.DB, startDate string, endDate string)([]JobStats, error){

	query := `SELECT b.job_type, b.status, count(b.*) AS count
            FROM (
              SELECT a.id,a.status,
                CASE WHEN array_length(a.types, 1) = 1 THEN a.types[1]
              ELSE 'Mixed'
           END AS job_type
           FROM (
               SELECT j.id, j.status, array_agg(DISTINCT t.name) AS types FROM jobs j
               JOIN job_steps s ON j.id = s.job_id
               JOIN job_types t ON s.job_type_id = t.id
               WHERE j.start_date >= ($1 :: DATE) AND j.start_date <= ($2 :: DATE) + INTERVAL '1 day'
               AND j.app_id != '1e8f719b-0452-4d39-a2f3-8714793ee3e6'
               AND j.status in ('Completed', 'Failed', 'Canceled')
               GROUP BY j.id
           ) AS a
          ) AS b
         GROUP BY b.job_type, b.status
         ORDER BY b.job_type;`

	rows, err := db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err;
	}

	defer rows.Close()

	var jobs []JobStats

	for rows.Next() {
		var job JobStats
		err := rows.Scan(&job.Category, &job.Status, &job.Count)
		if err != nil {
			return nil, err;
		}
		jobs = append(jobs, job)
		output := fmt.Sprintf("Total no.of %[1]v jobs %[2]v: %[3]v", job.Category, job.Status, job.Count)
		fmt.Println(output)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return jobs, nil


}