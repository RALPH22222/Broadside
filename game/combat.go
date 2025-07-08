package game

import (
	"fmt"
)

// Level constants
const (
	LevelFG = iota // Frigate (Easy)
	LevelDD        // Destroyer (Medium)
	LevelCB        // Cruiser (Hard)
	LevelBB        // Battleship (Expert)
)

var levelNames = []string{"Frigate", "Destroyer", "Cruiser", "Battleship"}

// Weapon struct
type Weapon struct {
	Name   string
	Damage float64 // multiplier (1.0 = 100%)
}

var (
	Cannon  = Weapon{"Cannon", 1.0}
	Torpedo = Weapon{"Torpedo", 1.1}
	Missile = Weapon{"Missile", 1.2}
	Railgun = Weapon{"Railgun", 1.3}
)

// Player struct
type Player struct {
	Shields int
	Streak  int
	Weapon  Weapon
	Score   int
}

// Enemy struct
type Enemy struct {
	HP    int
	Level int
}

// Add mainQuestions per level
var mainQuestions = []int{10, 15, 20, 25}

// Add damage reduction per level
var damageReduction = []float64{
	0.0, // Easy: 0% reduction
	0.2, // Medium: 20% reduction
	0.4, // Hard: 40% reduction
	0.6, // Extreme: 60% reduction
}

// Points per correct answer for each level
var basePoints = []int{
	10, // Easy
	20, // Medium
	30, // Hard
	40, // Extreme
}

// CalculatePoints returns points for a correct answer and any bonus points for remaining questions
// Parameters:
// - level: difficulty level
// - questionsLeft: number of unused questions remaining
// - bonusActive: whether bonus question was answered correctly
// - enemyDefeated: whether the enemy was defeated
func CalculatePoints(level int, questionsLeft int, bonusActive bool, enemyDefeated bool) (points int, bonus int) {
	// Base points for correct answer
	points = basePoints[level]

	// Calculate bonus points for remaining questions
	if enemyDefeated && bonusActive && questionsLeft > 0 {
		bonus = questionsLeft * basePoints[level]
	}

	return points, bonus
}

// NewCombat sets up a new combat scenario
func NewCombat(level int) (Player, Enemy) {
	return Player{Shields: 3, Streak: 0, Weapon: Cannon, Score: 0}, Enemy{HP: 100, Level: level}
}

// CalculateDamage computes the damage for a correct answer
// bonusActive: true if bonus question was answered correctly
func CalculateDamage(level int, weapon Weapon, bonusActive bool) int {
	// Calculate base damage needed to defeat enemy in exactly mainQuestions hits
	// Enemy has 100 HP, so we need: damage * mainQuestions = 100
	// This ensures player must answer all main questions correctly to win without bonus
	// Bonus question gives +10% damage if correct
	// Damage reduction is applied per level
	baseDamage := 100.0 / float64(mainQuestions[level])

	// Apply bonus damage if bonus question was answered correctly
	damage := baseDamage
	if bonusActive {
		damage *= 1.1 // +10% bonus damage
	}

	// Apply enemy damage reduction for this level
	damage *= (1.0 - damageReduction[level])

	return int(damage)
}

// ApplyAnswer updates player/enemy state based on answer
func ApplyAnswer(player *Player, enemy *Enemy, isCorrect, isBonus, isTimeout bool, questionsLeft int) {
	if isTimeout {
		player.Shields--
		player.Streak = 0
		fmt.Println("Timeout! Lost 1 shield.")
		return
	}

	if isCorrect {
		player.Streak++
		// Calculate combat damage
		dmg := CalculateDamage(enemy.Level, player.Weapon, isBonus)
		enemy.HP -= dmg

		// Calculate points
		points, bonus := CalculatePoints(enemy.Level, questionsLeft, isBonus, enemy.HP <= 0)
		player.Score += points + bonus

		if bonus > 0 {
			fmt.Printf("Correct! Dealt %d damage. Earned %d points + %d bonus points for remaining questions!\n",
				dmg, points, bonus)
		} else {
			fmt.Printf("Correct! Dealt %d damage. Earned %d points.\n", dmg, points)
		}

		if isBonus {
			switch enemy.Level {
			case LevelFG:
				player.Weapon = Torpedo
			case LevelDD:
				player.Weapon = Missile
			case LevelCB:
				player.Weapon = Railgun
			}
			fmt.Printf("Bonus! Upgraded weapon to %s.\n", player.Weapon.Name)
		}
	} else {
		player.Shields--
		player.Streak = 0
		fmt.Printf("Wrong! Lost 1 shield. (Shields left: %d)\n", player.Shields)
	}
}

// IsGameOver checks if the game is over
func IsGameOver(player Player, enemy Enemy) (bool, string) {
	if enemy.HP <= 0 {
		return true, "Victory! Enemy defeated."
	}
	if player.Shields <= 0 {
		return true, "Defeat! All shields down."
	}
	return false, ""
}

// Stub for asking a question (replace with real logic later)
func AskQuestion() (isCorrect, isBonus, isTimeout bool) {
	// For now, simulate always correct, not bonus, not timeout
	return true, false, false
}

// RunCombat runs the combat loop for a given level
func RunCombat(level int) {
	player, enemy := NewCombat(level)
	fmt.Printf("Starting battle: %s vs Enemy (HP: %d)\n", levelNames[level], enemy.HP)

	// Ask bonus question first
	fmt.Println("Bonus Question!")
	_, isBonus, _ := AskQuestion()
	if isBonus {
		ApplyAnswer(&player, &enemy, true, true, false, 0)
	}

	for {
		isCorrect, _, isTimeout := AskQuestion()
		ApplyAnswer(&player, &enemy, isCorrect, false, isTimeout, 0)
		if player.Shields <= 0 && enemy.HP > 0 {
			fmt.Println("Enemy fires back! Shields are down.")
		}
		gameOver, msg := IsGameOver(player, enemy)
		if gameOver {
			fmt.Println(msg)
			fmt.Printf("Final Score: %d\n", player.Score)
			break
		}
	}
}

// ProcessAnswer handles the result of a quiz answer and updates all combat state.
// Returns updated playerHP, playerShields, enemyHP, score, mainQDone, bonusActive, bonusAnswered, combatOver, and any other needed state.
type CombatState struct {
	PlayerHP      int
	PlayerShields int
	EnemyHP       int
	Score         int
	MainQDone     int
	BonusActive   bool
	BonusAnswered bool
	CombatOver    bool
	Rank          string
	ScorePercent  int
}

// ProcessAnswer processes a quiz answer and returns the updated combat state.
func ProcessAnswer(
	level int,
	playerHP int,
	playerShields int,
	enemyHP int,
	score int,
	mainQDone int,
	bonusActive bool,
	bonusAnswered bool,
	bonusQIndex int,
	currentQ int,
	isBonusQ bool,
	isCorrect bool,
	questionsLeft int,
) CombatState {
	// Copy input state
	newPlayerHP := playerHP
	newPlayerShields := playerShields
	newEnemyHP := enemyHP
	newScore := score
	newMainQDone := mainQDone
	newBonusActive := bonusActive
	newBonusAnswered := bonusAnswered
	combatOver := false

	// --- Handle Bonus Question ---
	if isBonusQ && !bonusAnswered {
		if isCorrect {
			// Bonus correct: activate bonus for rest of round
			newBonusActive = true
			newBonusAnswered = true
		} else {
			// Bonus wrong: no bonus
			newBonusAnswered = false
		}
		// No damage or points for bonus Q itself
		return CombatState{
			PlayerHP:      newPlayerHP,
			PlayerShields: newPlayerShields,
			EnemyHP:       newEnemyHP,
			Score:         newScore,
			MainQDone:     newMainQDone,
			BonusActive:   newBonusActive,
			BonusAnswered: newBonusAnswered,
			CombatOver:    false,
		}
	}

	// --- Handle Main Questions ---
	if isCorrect {
		// Player deals damage to enemy
		dmg := CalculateDamage(level, Cannon, newBonusActive)
		newEnemyHP -= dmg
		if newEnemyHP < 0 {
			newEnemyHP = 0
		}
		newMainQDone++
		// Award points for correct answer
		points, _ := CalculatePoints(level, 0, false, false)
		newScore += points
		// If enemy defeated and bonus was correct, award bonus points for remaining questions
		if newEnemyHP == 0 && newBonusActive && questionsLeft > 0 {
			newScore += questionsLeft * basePoints[level]
		}
		// Bonus stays active for the entire round once activated
	} else {
		// Wrong answer: lose shield or take damage
		if newPlayerShields > 0 {
			newPlayerShields--
		} else {
			// Enemy deals fixed damage per level
			enemyDmg := basePoints[level] // 10/20/30/40 per level
			newPlayerHP -= enemyDmg
			if newPlayerHP < 0 {
				newPlayerHP = 0
			}
		}
	}

	if newPlayerHP == 0 || newEnemyHP == 0 {
		combatOver = true
	}

	return CombatState{
		PlayerHP:      newPlayerHP,
		PlayerShields: newPlayerShields,
		EnemyHP:       newEnemyHP,
		Score:         newScore,
		MainQDone:     newMainQDone,
		BonusActive:   newBonusActive,
		BonusAnswered: newBonusAnswered,
		CombatOver:    combatOver,
	}
}

// GetMainQuestionsCount returns the number of main questions for a given level.
// This is used by the UI to determine how many questions to load for a combat session.
func GetMainQuestionsCount(level int) int {
	return mainQuestions[level]
}
