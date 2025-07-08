package game

// GameState tracks the current game session
type GameState struct {
	PlayerName      string
	Score           int
	Level           int
	Shields         int
	HP              int
	Streak          int
	Weapon          Weapon
	QuestsCompleted int
	WeaponBoosts    int
	CorrectAnswers  int
	TotalAnswers    int
	BonusSuccess    int
	BonusAttempts   int
}

// NewGameState creates a new game state for a player
func NewGameState(playerName string, level int) *GameState {
	return &GameState{
		PlayerName: playerName,
		Level:      level,
		Shields:    3,
		HP:         100,
		Weapon:     Weapon{"Cannon", 1.0},
	}
}

// AddScore adds to the player's score
func (gs *GameState) AddScore(points int) {
	gs.Score += points
}

// RecordAnswer updates stats for an answer
func (gs *GameState) RecordAnswer(correct bool) {
	gs.TotalAnswers++
	if correct {
		gs.CorrectAnswers++
		gs.Streak++
	} else {
		gs.Streak = 0
	}
}

// RecordBonus updates bonus question stats
func (gs *GameState) RecordBonus(success bool) {
	gs.BonusAttempts++
	if success {
		gs.BonusSuccess++
		gs.WeaponBoosts++
	}
}

// CompleteQuest increments quests completed
func (gs *GameState) CompleteQuest() {
	gs.QuestsCompleted++
}

// Accuracy returns the player's answer accuracy as a float
func (gs *GameState) Accuracy() float64 {
	if gs.TotalAnswers == 0 {
		return 0
	}
	return float64(gs.CorrectAnswers) / float64(gs.TotalAnswers)
}

// BonusSuccessRate returns the bonus question success rate
func (gs *GameState) BonusSuccessRate() float64 {
	if gs.BonusAttempts == 0 {
		return 0
	}
	return float64(gs.BonusSuccess) / float64(gs.BonusAttempts)
}

// LeaderboardEntry represents a leaderboard record
type LeaderboardEntry struct {
	PlayerName      string
	Score           int
	QuestsCompleted int
	WeaponBoosts    int
	Accuracy        float64
	BonusSuccess    float64
}

// NewLeaderboardEntry creates a leaderboard entry from a game state
func NewLeaderboardEntry(gs *GameState) LeaderboardEntry {
	return LeaderboardEntry{
		PlayerName:      gs.PlayerName,
		Score:           gs.Score,
		QuestsCompleted: gs.QuestsCompleted,
		WeaponBoosts:    gs.WeaponBoosts,
		Accuracy:        gs.Accuracy(),
		BonusSuccess:    gs.BonusSuccessRate(),
	}
}

// CalcRankAndPercent calculates the rank and score percent for the player.
func CalcRankAndPercent(score, mainQDone int, bonusAnswered bool, enemyHP, playerHP, level int) (rank string, percent int) {
	mainQuestions := []int{10, 15, 20, 25}
	basePoints := []int{10, 20, 30, 40}
	maxScore := mainQuestions[level] * basePoints[level]
	percent = 0
	if maxScore > 0 {
		percent = (score * 100) / maxScore
	}
	// If enemy defeated, force 100%
	if enemyHP == 0 {
		percent = 100
	}
	if playerHP == 0 || percent < 50 {
		rank = "Defeated"
	} else if percent == 100 {
		rank = "S+"
	} else if percent >= 90 {
		rank = "Gold"
	} else if percent >= 70 {
		rank = "Silver"
	} else if percent >= 50 {
		rank = "Bronze"
	}
	return
}

// CombatInitState holds all initial values for a new combat session.
type CombatInitState struct {
	Level            int
	PlayerMaxHP      int
	PlayerHP         int
	PlayerMaxShields int
	PlayerShields    int
	EnemyMaxHP       int
	EnemyHP          int
	BonusActive      bool
	BonusAnswered    bool
	BonusQDone       bool
	MainQDone        int
	CombatOver       bool
}

// InitCombatState initializes all combat state for a new session based on difficulty.
// This keeps UI and game logic cleanly separated for maintainability.
func InitCombatState(difficulty string) CombatInitState {
	levelMap := map[string]int{"Easy": 0, "Medium": 1, "Hard": 2, "Extreme": 3}
	level := levelMap[difficulty]
	return CombatInitState{
		Level:            level,
		PlayerMaxHP:      100,
		PlayerHP:         100,
		PlayerMaxShields: 3,
		PlayerShields:    3,
		EnemyMaxHP:       100,
		EnemyHP:          100,
		BonusActive:      false,
		BonusAnswered:    false,
		BonusQDone:       false,
		MainQDone:        0,
		CombatOver:       false,
	}
}
