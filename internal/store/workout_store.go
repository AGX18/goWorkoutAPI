package store

import "database/sql"

type Workout struct {
	ID              string         `json:"id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	DurationMinutes int            `json:"duration_minutes"` // Duration in minutes
	CaloriesBurned  int            `json:"calories_burned"`  // Calories burned during the workout
	Entries         []WorkoutEntry `json:"entries"`          // List of exercise entries in the workout
}

type WorkoutEntry struct {
	ID              string   `json:"id"`
	ExerciseName    string   `json:"exercise_name"`
	Sets            int      `json:"sets"`
	Reps            *int     `json:"reps"`
	DurationSeconds *int     `json:"duration_seconds"` // Duration in seconds, if applicable
	Weight          *float64 `json:"weight"`
	Notes           *string  `json:"notes"`       // Additional notes for the entry
	OrderIndex      int      `json:"order_index"` // Order of the entry in the workout
}

type PostgresWorkoutStore struct {
	db *sql.DB
}

func NewPostgresWorkoutStore(db *sql.DB) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

type WorkoutStore interface {
	CreateWorkout(workout *Workout) (*Workout, error)
	GetWorkoutByID(id int64) (*Workout, error)
}

func (pg *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query :=
		`
  INSERT INTO workouts (title, description, duration_minutes, calories_burned)
  VALUES ($1, $2, $3, $4)
  RETURNING id 
  `

	err = tx.QueryRow(query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)
	if err != nil {
		return nil, err
	}

	for _, entry := range workout.Entries {
		query := `
		INSERT INTO workout_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
    	`
		err = tx.QueryRow(query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)
		if err != nil {
			return nil, err
		}

	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {

	workout := &Workout{}

	return workout, nil

}
