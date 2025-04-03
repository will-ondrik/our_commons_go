package database

import (
	"context"
	"etl_our_commons/dtos"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InsertAll inserts MPs, MP Terms, Expenses, and related logs in batch
func (db *DB) InsertAll(pool *pgxpool.Pool, mps []*dtos.Mp) error {
	// Insert MPs
	fmt.Println("[INSERTING]: Mps...")
	mpIds, err := db.InsertMPs(pool, mps)
	if err != nil {
		return fmt.Errorf("failed to insert MPs: %w", err)
	}

	// Insert MP Terms
	fmt.Println("[INSERTING]: Mp Terms...")
	mpTermIds, err := db.InsertMPTerms(pool, mps, mpIds)
	if err != nil {
		return fmt.Errorf("failed to insert MP terms: %w", err)
	}
	fmt.Println("Completed successfully.")

	// Insert Expenses (Salaries, Contract, Hospitality, Travel)
	fmt.Println("[INSERTING]: Expenses...")
	expenseIdsMap := make(map[string]map[string]int)
	for _, mp := range mps {
		ids, err := db.InsertExpenses(pool, mp, mpTermIds)
		if err != nil {
			return fmt.Errorf("failed to insert expenses for %s: %w", mp.MpName.FirstName, err)
		}
		expenseIdsMap[mp.MpName.FirstName+"_"+mp.MpName.LastName] = ids
	}
	fmt.Println("Completed successfully.")

	// Insert Hospitality Events
	fmt.Println("[INSERTING]: Hospitality events...")
	for _, mp := range mps {
		mpKey := mp.MpName.FirstName + "_" + mp.MpName.LastName
		expenseIds, exists := expenseIdsMap[mpKey]
		if !exists {
			return fmt.Errorf("expense IDs not found for %s", mpKey)
		}
		_, err = db.InsertHospitalityEvents(pool, mp, expenseIds)
		if err != nil {
			return fmt.Errorf("failed to insert hospitality events for %s: %w", mp.MpName.FirstName, err)
		}
	}
	fmt.Println("Completed successfully.")

	// Insert Travel Events
	fmt.Println("[INSERTING]: Travel events...")
	travelEventIds := make(map[string]map[string]int)
	for _, mp := range mps {
		mpKey := mp.MpName.FirstName + "_" + mp.MpName.LastName
		expenseIds, exists := expenseIdsMap[mpKey]
		if !exists {
			return fmt.Errorf("expense IDs not found for %s", mpKey)
		}
		eventIds, err := db.InsertTravelEvents(pool, mp, expenseIds)
		if err != nil {
			return fmt.Errorf("failed to insert travel events for %s: %w", mp.MpName.FirstName, err)
		}
		travelEventIds[mpKey] = eventIds
	}
	fmt.Println("Completed successfully.")

	// Insert Travellers
	fmt.Println("[INSERTING]: Travellers...")
	travellerIds := make(map[string]int)
	for _, mp := range mps {
		ids, err := db.InsertTravellers(pool, mp)
		if err != nil {
			return fmt.Errorf("failed to insert travellers for %s: %w", mp.MpName.FirstName, err)
		}
		for k, v := range ids {
			travellerIds[k] = v
		}
	}
	fmt.Println("Completed successfully.")


	// Insert Traveller Logs
	fmt.Println("[INSERTING]: Traveller logs...")
	travellerLogIds := make(map[string]map[string]int)
	for _, mp := range mps {
		mpKey := mp.MpName.FirstName + "_" + mp.MpName.LastName
		expenseIds, exists := expenseIdsMap[mpKey]
		if !exists {
			return fmt.Errorf("expense IDs not found for %s", mpKey)
		}
		logIds, err := db.InsertTravellerLogs(pool, mp, expenseIds)
		if err != nil {
			return fmt.Errorf("failed to insert traveller logs for %s: %w", mp.MpName.FirstName, err)
		}
		travellerLogIds[mpKey] = logIds
	}
	fmt.Println("Completed successfully.")

	// Insert Travel Data and Flight Points
	fmt.Println("[INSERTING]: Travel data and flight points...")
	for _, mp := range mps {
		mpKey := mp.MpName.FirstName + "_" + mp.MpName.LastName
		if eventIds, exists := travelEventIds[mpKey]; exists {
			travelDataIds, err := db.InsertTravelData(pool, mp, eventIds)
			if err != nil {
				return fmt.Errorf("failed to insert travel data for %s: %w", mp.MpName.FirstName, err)
			}

			err = db.InsertFlightPoints(pool, mp, travelDataIds)
			if err != nil {
				return fmt.Errorf("failed to insert flight points for %s: %w", mp.MpName.FirstName, err)
			}
		}
	}
	fmt.Println("Completed successfully.")

	// Insert Quarterly Reports
	fmt.Println("[INSERTING]: Quarterly reports...")
	if len(mps) > 0 {
		// We only need to insert quarterly reports once, not for each MP
		err = db.InsertQuarterlyReports(pool, mps[0])
		if err != nil {
			return fmt.Errorf("failed to insert quarterly reports: %w", err)
		}
	}
	fmt.Println("Completed successfully.")

	fmt.Println("All records inserted successfully!")
	return nil
}


// InsertExpenses inserts expenses per MP
func (db *DB) InsertExpenses(pool *pgxpool.Pool, mp *dtos.Mp, mpTermIds map[string]int) (map[string]int, error) {
	expenseIds := make(map[string]int)
	ctx := context.Background()

	mpKey := mp.Constituency + "_" + mp.Caucus + "_" + fmt.Sprint(mp.FiscalYear) + "_Q" + fmt.Sprint(mp.FiscalQuarter)
	mpTermId, exists := mpTermIds[mpKey]
	if !exists {
		return nil, fmt.Errorf("MP Term ID not found for %s %s", mp.MpName.FirstName, mp.MpName.LastName)
	}

	// Salary Expenses
	if mp.Expenses.Totals.SalariesCost > 0 {
		var salaryId int
		err := pool.QueryRow(ctx, `
			INSERT INTO salaryexpenses (mpterm, amount)
			VALUES ($1, $2)
			ON CONFLICT (mpterm) DO NOTHING
			RETURNING id;
		`, mpTermId, mp.Expenses.Totals.SalariesCost).Scan(&salaryId)

		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		} else if err == pgx.ErrNoRows {
			err = pool.QueryRow(ctx, `SELECT id FROM salaryexpenses WHERE mpterm = $1;`, mpTermId).Scan(&salaryId)
			if err != nil {
				return nil, err
			}
		}

		expenseIds["salary"] = salaryId
	}

	// Contract Expenses
	for i, contract := range mp.Expenses.ContractExpenses {
		var contractId int
		err := pool.QueryRow(ctx, `
			INSERT INTO contractexpenses (mpterm, amount, url)
			VALUES ($1, $2, $3)
			ON CONFLICT (mpterm) DO UPDATE SET url = EXCLUDED.url
			RETURNING id;
		`, mpTermId, contract.Total, mp.Expenses.ContractExpensesUrl).Scan(&contractId)

		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		} else if err == pgx.ErrNoRows {
			err = pool.QueryRow(ctx, `SELECT id FROM contractexpenses WHERE mpterm = $1;`, mpTermId).Scan(&contractId)
			if err != nil {
				return nil, err
			}
		}

		expenseIds["contract_"+strconv.Itoa(i)] = contractId

		_, err = pool.Exec(ctx, `
			INSERT INTO contractexpenselogs (expenseid, supplier, description, expensedate, amount)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING;
		`, contractId, contract.Supplier, contract.Description, contract.Date, contract.Total)

		if err != nil {
			return nil, err
		}
	}

	// Hospitality Expenses
	if len(mp.Expenses.HospitalityExpenses) > 0 {
		var hospitalityId int
		err := pool.QueryRow(ctx, `
			INSERT INTO hospitalityexpenses (mpterm, amount, url)
			VALUES ($1, $2, $3)
			ON CONFLICT (mpterm) DO UPDATE SET url = EXCLUDED.url
			RETURNING id;
		`, mpTermId, mp.Expenses.HospitalityExpenses[0].TotalCost, mp.Expenses.HospitalityExpensesUrl).Scan(&hospitalityId)

		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		} else if err == pgx.ErrNoRows {
			err = pool.QueryRow(ctx, `SELECT id FROM hospitalityexpenses WHERE mpterm = $1;`, mpTermId).Scan(&hospitalityId)
			if err != nil {
				return nil, err
			}
		}

		expenseIds["hospitality"] = hospitalityId
	}

	// Travel Expenses
	if len(mp.Expenses.TravelExpenses) > 0 {
		var travelId int
		err := pool.QueryRow(ctx, `
			INSERT INTO travelexpenses (mpterm, amount, url)
			VALUES ($1, $2, $3)
			ON CONFLICT (mpterm) DO UPDATE SET url = EXCLUDED.url
			RETURNING id;
		`, mpTermId, mp.Expenses.TravelExpenses[0].TravelCosts.Total, mp.Expenses.TravelExpensesUrl).Scan(&travelId)

		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		} else if err == pgx.ErrNoRows {
			err = pool.QueryRow(ctx, `SELECT id FROM travelexpenses WHERE mpterm = $1;`, mpTermId).Scan(&travelId)
			if err != nil {
				return nil, err
			}
		}

		expenseIds["travel"] = travelId
	}

	return expenseIds, nil
}


// InsertMPs inserts MPs into the database and returns a mapping of MP names to their IDs
func (db *DB) InsertMPs(pool *pgxpool.Pool, mps []*dtos.Mp) (map[string]int, error) {
	mpIds := make(map[string]int)
	ctx := context.Background()

	// Process MPs one by one instead of in a batch for better error handling
	for _, mp := range mps {
		// Check for empty firstName
		if mp.MpName.FirstName == "" {
			return nil, fmt.Errorf("MP has empty firstName: lastName=%s, constituency=%s, caucus=%s", 
				mp.MpName.LastName, mp.Constituency, mp.Caucus)
		}

		var id int
		err := pool.QueryRow(ctx, `
			INSERT INTO mp (firstname, lastname, constituency, caucus)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (firstname, lastname) DO UPDATE 
			SET constituency = $3, caucus = $4
			RETURNING id;
		`, mp.MpName.FirstName, mp.MpName.LastName, mp.Constituency, mp.Caucus).Scan(&id)
		
		if err != nil {
			return nil, fmt.Errorf("failed to insert MP %s %s: %w", 
				mp.MpName.FirstName, mp.MpName.LastName, err)
		}
		
		mpIds[mp.MpName.FirstName+"_"+mp.MpName.LastName] = id
	}

	return mpIds, nil
}


// InsertMPTerms inserts MP terms into the database and returns a mapping of MP term identifiers to their IDs
func (db *DB) InsertMPTerms(pool *pgxpool.Pool, mps []*dtos.Mp, mpIds map[string]int) (map[string]int, error) {
	mpTermIds := make(map[string]int)
	batch := &pgx.Batch{}
	ctx := context.Background()

	for _, mp := range mps {
		mpId, exists := mpIds[mp.MpName.FirstName+"_"+mp.MpName.LastName]
		if !exists {
			return nil, fmt.Errorf("MP ID not found for %s %s", mp.MpName.FirstName, mp.MpName.LastName)
		}

		batch.Queue(`
			INSERT INTO mpterm (mpid, startdate, enddate, fiscalyear, fiscalquarter, constituency, caucus, url) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (mpid, fiscalyear, fiscalquarter) DO UPDATE SET url = EXCLUDED.url 
			RETURNING id;
		`, mpId, mp.Years.StartDate, mp.Years.EndDate, mp.FiscalYear, mp.FiscalQuarter, mp.Constituency, mp.Caucus, mp.Url)
	}

	batchResults := pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	for _, mp := range mps {
		var id int
		err := batchResults.QueryRow().Scan(&id)
		if err != nil {
			if err == pgx.ErrNoRows {
				mpId := mpIds[mp.MpName.FirstName+"_"+mp.MpName.LastName]
				err = pool.QueryRow(ctx, `
					SELECT id FROM mpterm WHERE mpid = $1 AND fiscalyear = $2 AND fiscalquarter = $3;
				`, mpId, mp.FiscalYear, mp.FiscalQuarter).Scan(&id)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		mpTermIds[mp.Constituency+"_"+mp.Caucus+"_"+fmt.Sprint(mp.FiscalYear)+"_Q"+fmt.Sprint(mp.FiscalQuarter)] = id
	}

	return mpTermIds, nil
}

// InsertHospitalityEvents inserts hospitality events for each MP into the database
// Uses a transaction to ensure all operations are performed on the same connection
func (db *DB) InsertHospitalityEvents(pool *pgxpool.Pool, mp *dtos.Mp, expenseIds map[string]int) (map[string]int, error) {
	eventIds := make(map[string]int)
	ctx := context.Background()

	hospitalityExpenseId, exists := expenseIds["hospitality"]
	if !exists {
		return eventIds, nil
	}

	// Start a transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Default date to use when date strings are empty
	const defaultDate = "1970-01-01"

	// Insert hospitality events
	for _, expense := range mp.Expenses.HospitalityExpenses {
		// Check for empty date and use default if needed
		expenseDate := expense.Date
		if expenseDate == "" {
			fmt.Printf("Warning: Empty date for hospitality event, using default date\n")
			expenseDate = defaultDate
		}

		var id int
		err := tx.QueryRow(ctx, `
			INSERT INTO hospitalityevents (expenseid, expensedate, location, purpose, amount)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (expenseid, expensedate) DO NOTHING
			RETURNING id;
		`, hospitalityExpenseId, expenseDate, expense.Location, expense.Purpose, expense.TotalCost).Scan(&id)

		if err != nil && err != pgx.ErrNoRows {
			return nil, fmt.Errorf("failed to insert hospitality event: %w", err)
		}

		if err == pgx.ErrNoRows {
			// If no rows were returned, get the existing ID
			err = tx.QueryRow(ctx, `
				SELECT id FROM hospitalityevents 
				WHERE expenseid = $1 AND expensedate = $2;
			`, hospitalityExpenseId, expenseDate).Scan(&id)
			if err != nil {
				return nil, fmt.Errorf("failed to get existing hospitality event ID: %w", err)
			}
		}

		eventIds[expense.Date+"_"+expense.Location] = id

		// Insert expense logs for this event
		for _, log := range expense.Event.ExpenseLogs {
			_, err = tx.Exec(ctx, `
				INSERT INTO hospitalityexpenselogs (eventid, claimnumber, supplier, amount)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING;
			`, id, log.Claim, log.Supplier, log.Cost)

			if err != nil {
				return nil, fmt.Errorf("failed to insert hospitality expense log: %w", err)
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return eventIds, nil
}

// InsertTravellerLogs inserts traveller logs related to travel expenses
func (db *DB) InsertTravellerLogs(pool *pgxpool.Pool, mp *dtos.Mp, expenseIds map[string]int) (map[string]int, error) {
	logIds := make(map[string]int)
	ctx := context.Background()

	travelExpenseId, exists := expenseIds["travel"]
	if !exists {
		return logIds, nil
	}

	// Get travel events for this expense
	rows, err := pool.Query(ctx, `
		SELECT id, claimnumber FROM travelevents WHERE expenseid = $1
	`, travelExpenseId)
	if err != nil {
		return nil, fmt.Errorf("failed to query travel events: %w", err)
	}
	defer rows.Close()

	// Map claim numbers to event IDs
	eventMap := make(map[string]int)
	for rows.Next() {
		var id int
		var claim string
		if err := rows.Scan(&id, &claim); err != nil {
			return nil, fmt.Errorf("failed to scan travel event: %w", err)
		}
		eventMap[claim] = id
	}

	// Get traveller IDs
	travellerIds := make(map[string]int)
	for _, travel := range mp.Expenses.TravelExpenses {
		for _, traveller := range travel.TravelLogs {
			key := traveller.Name.FirstName + "_" + traveller.Name.LastName
			if _, exists := travellerIds[key]; exists {
				continue
			}

			var id int
			err := pool.QueryRow(ctx, `
				SELECT id FROM traveller WHERE firstname = $1 AND lastname = $2 AND type = $3
			`, traveller.Name.FirstName, traveller.Name.LastName, traveller.Type).Scan(&id)

			if err != nil {
				return nil, fmt.Errorf("failed to get traveller ID: %w", err)
			}

			travellerIds[key] = id
		}
	}

	// Now insert traveller logs with correct eventid and travellerid
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	for _, travel := range mp.Expenses.TravelExpenses {
		eventId, exists := eventMap[travel.Claim]
		if !exists {
			// If we don't have an event ID for this claim, skip it
			continue 
		}

		for _, traveller := range travel.TravelLogs {
			travellerKey := traveller.Name.FirstName + "_" + traveller.Name.LastName
			travellerId, exists := travellerIds[travellerKey]
			if !exists {
				// If we don't have a traveller ID, skip it
				continue 
			}

			var id int
			var travelDate interface{} = traveller.Date
			
			// Convert empty string to nil (NULL in database)
			if traveller.Date == "" {
				travelDate = nil
				fmt.Printf("Using NULL for travel date for %s %s\n", traveller.Name.FirstName, traveller.Name.LastName)
			}

			err := tx.QueryRow(ctx, `
				INSERT INTO travellerlogs (eventid, travellerid, traveldate, purpose, departurecity, destinationcity, transportationmode)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (eventid, travellerid, traveldate) DO NOTHING
				RETURNING id;
			`, eventId, travellerId, travelDate, traveller.Purpose, traveller.DepartureCity, traveller.DestinationCity, traveller.TransportationMode).Scan(&id)

			if err != nil && err != pgx.ErrNoRows {
				return nil, fmt.Errorf("failed to insert traveller log: %w", err)
			}

			if err == pgx.ErrNoRows {

				// If no rows were returned, get the existing ID
				err = tx.QueryRow(ctx, `
					SELECT id FROM travellerlogs 
					WHERE eventid = $1 AND travellerid = $2 AND traveldate = $3;
				`, eventId, travellerId, traveller.Date).Scan(&id)
				if err != nil {
					return nil, fmt.Errorf("failed to get existing traveller log ID: %w", err)
				}
			}

			logIds[travellerKey] = id
		}
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return logIds, nil
}

// InsertTravelExpenseLogs inserts logs for travel expenses
func (db *DB) InsertTravelExpenseLogs(pool *pgxpool.Pool, mp *dtos.Mp, expenseIds map[string]int) error {
	batch := &pgx.Batch{}
	ctx := context.Background()

	travelExpenseId, exists := expenseIds["travel"]
	if !exists {
		return nil
	}

	for _, travel := range mp.Expenses.TravelExpenses {
		batch.Queue(`
			INSERT INTO travelexpenselogs (expenseid, transportationamount, accomodationamount, mealsandincidentalsamount, totalamount)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (expenseid) DO NOTHING;
		`, travelExpenseId, travel.TravelCosts.Transportation, travel.TravelCosts.Accomodation, travel.TravelCosts.MealsAndIncidentals, travel.TravelCosts.Total)
	}

	batchResults := pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	for range mp.Expenses.TravelExpenses {
		if err := batchResults.QueryRow().Scan(); err != nil && err != pgx.ErrNoRows {
			return err
		}
	}

	return nil
}

// InsertTravelData inserts travel-related data including distance and emissions
func (db *DB) InsertTravelData(pool *pgxpool.Pool, mp *dtos.Mp, logIds map[string]int) (map[string]int, error) {
	ctx := context.Background()
	travelDataIds := make(map[string]int)

	for _, travel := range mp.Expenses.TravelExpenses {
		eventId, exists := logIds[travel.Claim]
		if !exists {
			continue
		}

		var travelDataId int
		err := pool.QueryRow(ctx, `
			INSERT INTO traveldata (eventid, traveldistance, travelunit, emissions, emissionsunit)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (eventid) DO NOTHING
			RETURNING id;
		`, eventId, travel.TravelCosts.Total, "km", travel.TravelCosts.Total*0.2, "kg CO2").Scan(&travelDataId)

		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		} else if err == pgx.ErrNoRows {
			err = pool.QueryRow(ctx, `
				SELECT id FROM traveldata WHERE eventid = $1;
			`, eventId).Scan(&travelDataId)
			if err != nil {
				return nil, err
			}
		}

		travelDataIds[travel.Claim] = travelDataId
	}

	return travelDataIds, nil
}

// InsertFlightPoints inserts flight-related points data (FIXED: use traveldataid, not eventid)
func (db *DB) InsertFlightPoints(pool *pgxpool.Pool, mp *dtos.Mp, travelDataIds map[string]int) error {
	batch := &pgx.Batch{}
	ctx := context.Background()

	for _, travel := range mp.Expenses.TravelExpenses {
		travelDataId, exists := travelDataIds[travel.Claim]
		if !exists {
			continue
		}

		batch.Queue(`
			INSERT INTO flightpoints (traveldataid, regular, special, usa)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (traveldataid) DO NOTHING;
		`, travelDataId, travel.FlightPoints.Regular, travel.FlightPoints.Special, travel.FlightPoints.USA)
	}

	batchResults := pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	for range mp.Expenses.TravelExpenses {
		if err := batchResults.QueryRow().Scan(); err != nil && err != pgx.ErrNoRows {
			return err
		}
	}

	return nil
}
// InsertTravellers inserts traveller records from the TravelLogs data
// Uses a transaction to ensure all operations are performed on the same connection
func (db *DB) InsertTravellers(pool *pgxpool.Pool, mp *dtos.Mp) (map[string]int, error) {
	travellerIds := make(map[string]int)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	for _, travel := range mp.Expenses.TravelExpenses {
		for _, traveller := range travel.TravelLogs {
			key := traveller.Name.FirstName + "_" + traveller.Name.LastName
			if _, exists := travellerIds[key]; exists {
				continue
			}

			var id int
			err := tx.QueryRow(ctx, `
				INSERT INTO traveller (firstname, lastname, type)
				VALUES ($1, $2, $3)
				ON CONFLICT (firstname, lastname, type)
				DO UPDATE SET firstname = EXCLUDED.firstname
				RETURNING id;
			`, traveller.Name.FirstName, traveller.Name.LastName, traveller.Type).Scan(&id)

			if err != nil {
				return nil, fmt.Errorf("failed to insert traveller: %w", err)
			}

			travellerIds[key] = id
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return travellerIds, nil
}

// InsertQuarterlyReports inserts quarterly report data
func (db *DB) InsertQuarterlyReports(pool *pgxpool.Pool, mp *dtos.Mp) error {
	batch := &pgx.Batch{}
	ctx := context.Background()

	// Insert the quarterly report for this MP
	batch.Queue(`
		INSERT INTO quarterlyreports (startdate, enddate, fiscalyear, fiscalquarter, href)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (href) DO NOTHING;
	`, mp.Years.StartDate, mp.Years.EndDate, mp.FiscalYear, mp.FiscalQuarter, mp.Url)

	batchResults := pool.SendBatch(ctx, batch)
	defer batchResults.Close()

	// Check for errors
	if err := batchResults.QueryRow().Scan(); err != nil && err != pgx.ErrNoRows {
		return err
	}

	return nil
}

// InsertTravelEvents inserts travel events for each MP into the database
// Uses a transaction to ensure all operations are performed on the same connection
func (db *DB) InsertTravelEvents(pool *pgxpool.Pool, mp *dtos.Mp, expenseIds map[string]int) (map[string]int, error) {
	eventIds := make(map[string]int)
	ctx := context.Background()

	travelExpenseId, exists := expenseIds["travel"]
	if !exists {
		return eventIds, nil
	}

	// Start a transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Default date to use when date strings are empty
	const defaultDate = "1970-01-01"

	// Insert travel events
	for _, travel := range mp.Expenses.TravelExpenses {
		var startDate, endDate interface{} = travel.Dates.StartDate, travel.Dates.EndDate

		// Convert "Not Provided" to nil (NULL in database)
		if travel.Dates.StartDate == "Not Provided" || travel.Dates.StartDate == "" {
			startDate = nil
			fmt.Printf("Using NULL for start date for claim %s\n", travel.Claim)
		}
		if travel.Dates.EndDate == "Not Provided" || travel.Dates.EndDate == "" {
			endDate = nil
			fmt.Printf("Using NULL for end date for claim %s\n", travel.Claim)
		}

		var id int
		err := tx.QueryRow(ctx, `
			INSERT INTO travelevents (expenseid, claimnumber, startdate, enddate)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (expenseid, claimnumber) DO NOTHING
			RETURNING id;
		`, travelExpenseId, travel.Claim, startDate, endDate).Scan(&id)

		if err != nil {
			if err == pgx.ErrNoRows {
				
				// Already exists, get existing ID
				err = tx.QueryRow(ctx, `
					SELECT id FROM travelevents 
					WHERE expenseid = $1 AND claimnumber = $2;
				`, travelExpenseId, travel.Claim).Scan(&id)
				if err != nil {
					return nil, fmt.Errorf("failed to fetch existing travel event ID: %w", err)
				}
			} else {
				return nil, fmt.Errorf("failed to insert travel event: %w", err)
			}
		}

		eventIds[travel.Claim] = id

		// Insert travel expense log for the event
		_, err = tx.Exec(ctx, `
			INSERT INTO travelexpenselogs (eventid, transportationamount, accomodationamount, mealsandincidentalsamount, totalamount)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING;
		`, id, travel.TravelCosts.Transportation, travel.TravelCosts.Accomodation, travel.TravelCosts.MealsAndIncidentals, travel.TravelCosts.Total)

		if err != nil {
			return nil, fmt.Errorf("failed to insert travel expense log: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return eventIds, nil
}
