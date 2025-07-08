package game

import (
	"math/rand"
	"strings"
	"time"
)

// Question represents a quiz question
type Question struct {
	Text       string
	Choices    []string
	Answer     string // correct answer (case-insensitive)
	Subject    string
	Difficulty string // e.g., "Easy", "Medium", "Hard", "Expert"
}

// Quiz holds all questions
type Quiz struct {
	Questions []Question
}

var sampleQuestions = []Question{
	// Math(easy)
	{Text: "What is 2 + 3?", Choices: []string{"4", "5", "6"}, Answer: "5", Subject: "Math", Difficulty: "Easy"},
	{Text: "Which number comes after 7?", Choices: []string{"6", "8", "9"}, Answer: "8", Subject: "Math", Difficulty: "Easy"},
	{Text: "What is 5 - 2?", Choices: []string{"3", "2", "4"}, Answer: "3", Subject: "Math", Difficulty: "Easy"},
	{Text: "How many sides does a triangle have?", Choices: []string{"3", "4", "5"}, Answer: "3", Subject: "Math", Difficulty: "Easy"},
	{Text: "What is the number before 10?", Choices: []string{"9", "8", "11"}, Answer: "9", Subject: "Math", Difficulty: "Easy"},
	{Text: "Which is more: 6 or 9?", Choices: []string{"6", "9", "They are equal"}, Answer: "9", Subject: "Math", Difficulty: "Easy"},
	{Text: "What shape is a wheel?", Choices: []string{"Square", "Circle", "Triangle"}, Answer: "Circle", Subject: "Math", Difficulty: "Easy"},
	{Text: "What is 1 + 1?", Choices: []string{"1", "2", "3"}, Answer: "2", Subject: "Math", Difficulty: "Easy"},
	{Text: "How many legs do two dogs have?", Choices: []string{"4", "8", "6"}, Answer: "8", Subject: "Math", Difficulty: "Easy"},
	{Text: "Which of these is the smallest number?", Choices: []string{"3", "1", "2"}, Answer: "1", Subject: "Math", Difficulty: "Easy"},
	{Text: "What time is it if the clock shows 12 and 0 minutes?", Choices: []string{"12 o'clock", "1 o'clock", "11 o'clock"}, Answer: "12 o'clock", Subject: "Math", Difficulty: "Easy"},
	// Math(medium)
	{Text: "What is 6 + 7?", Choices: []string{"13", "12", "14"}, Answer: "13", Subject: "Math", Difficulty: "Medium"},
	{Text: "What is 10 - 4?", Choices: []string{"5", "6", "7"}, Answer: "6", Subject: "Math", Difficulty: "Medium"},
	{Text: "Which number is greater: 15 or 12?", Choices: []string{"12", "15", "They are equal"}, Answer: "15", Subject: "Math", Difficulty: "Medium"},
	{Text: "What is the next number in the pattern: 2, 4, 6, ?", Choices: []string{"8", "7", "10"}, Answer: "8", Subject: "Math", Difficulty: "Medium"},
	{Text: "Which shape has 4 equal sides?", Choices: []string{"Circle", "Triangle", "Square"}, Answer: "Square", Subject: "Math", Difficulty: "Medium"},
	{Text: "What is 3 + 9?", Choices: []string{"11", "12", "13"}, Answer: "12", Subject: "Math", Difficulty: "Medium"},
	{Text: "What is 14 - 5?", Choices: []string{"9", "10", "8"}, Answer: "9", Subject: "Math", Difficulty: "Medium"},
	{Text: "How many tens are there in 30?", Choices: []string{"2", "3", "4"}, Answer: "3", Subject: "Math", Difficulty: "Medium"},
	{Text: "Which is the smallest: 17, 13, or 15?", Choices: []string{"17", "13", "15"}, Answer: "13", Subject: "Math", Difficulty: "Medium"},
	{Text: "How many sides does a rectangle have?", Choices: []string{"3", "4", "5"}, Answer: "4", Subject: "Math", Difficulty: "Medium"},
	{Text: "What number comes next: 5, 10, 15, ?", Choices: []string{"20", "25", "30"}, Answer: "20", Subject: "Math", Difficulty: "Medium"},
	{Text: "If you have 4 apples and get 3 more, how many apples do you have?", Choices: []string{"6", "7", "8"}, Answer: "7", Subject: "Math", Difficulty: "Medium"},
	{Text: "What is 20 - 8?", Choices: []string{"12", "11", "13"}, Answer: "12", Subject: "Math", Difficulty: "Medium"},
	{Text: "How many legs do 3 cats have?", Choices: []string{"8", "10", "12"}, Answer: "12", Subject: "Math", Difficulty: "Medium"},
	{Text: "What is 11 + 4?", Choices: []string{"15", "14", "13"}, Answer: "15", Subject: "Math", Difficulty: "Medium"},
	{Text: "Which is more: 7 tens or 60?", Choices: []string{"60", "70", "They are equal"}, Answer: "70", Subject: "Math", Difficulty: "Medium"},
	// Math(hard)
	{Text: "What is 9 + 6?", Choices: []string{"14", "15", "16"}, Answer: "15", Subject: "Math", Difficulty: "Hard"},
	{Text: "What number is missing? 2, 4, __, 8", Choices: []string{"5", "6", "7"}, Answer: "6", Subject: "Math", Difficulty: "Hard"},
	{Text: "Which number is in the tens place in 47?", Choices: []string{"4", "7", "0"}, Answer: "4", Subject: "Math", Difficulty: "Hard"},
	{Text: "Tom has 3 red balls and 4 blue balls. How many balls does he have in total?", Choices: []string{"6", "7", "8"}, Answer: "7", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is 10 - 7 + 2?", Choices: []string{"5", "4", "3"}, Answer: "5", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is the largest number? 21, 12, or 19?", Choices: []string{"12", "19", "21"}, Answer: "21", Subject: "Math", Difficulty: "Hard"},
	{Text: "Which is an even number?", Choices: []string{"5", "7", "8"}, Answer: "8", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is 3 + 3 + 3?", Choices: []string{"9", "6", "8"}, Answer: "9", Subject: "Math", Difficulty: "Hard"},
	{Text: "Which number is 1 more than 99?", Choices: []string{"100", "98", "101"}, Answer: "100", Subject: "Math", Difficulty: "Hard"},
	{Text: "How many sides does a rectangle have?", Choices: []string{"3", "4", "5"}, Answer: "4", Subject: "Math", Difficulty: "Hard"},
	{Text: "If you count by 5s starting from 5, what comes after 15?", Choices: []string{"20", "25", "10"}, Answer: "20", Subject: "Math", Difficulty: "Hard"},
	{Text: "You have 2 boxes. One has 6 apples and the other has 4. How many apples in total?", Choices: []string{"10", "9", "11"}, Answer: "10", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is double of 7?", Choices: []string{"13", "14", "15"}, Answer: "14", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is half of 10?", Choices: []string{"4", "5", "6"}, Answer: "5", Subject: "Math", Difficulty: "Hard"},
	{Text: "Which group has more: 3 birds or 5 birds?", Choices: []string{"3 birds", "5 birds", "They are equal"}, Answer: "5 birds", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is 12 - 4 - 3?", Choices: []string{"6", "5", "4"}, Answer: "5", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is the smallest two-digit number?", Choices: []string{"10", "11", "12"}, Answer: "10", Subject: "Math", Difficulty: "Hard"},
	{Text: "Which shape has 4 equal sides?", Choices: []string{"Rectangle", "Square", "Triangle"}, Answer: "Square", Subject: "Math", Difficulty: "Hard"},
	{Text: "You have 5 pencils and give away 2. How many do you have left?", Choices: []string{"2", "3", "4"}, Answer: "3", Subject: "Math", Difficulty: "Hard"},
	{Text: "What is 8 + 2 - 5?", Choices: []string{"4", "5", "6"}, Answer: "5", Subject: "Math", Difficulty: "Hard"},
	{Text: "Which number comes next? 11, 13, 15, __", Choices: []string{"17", "18", "16"}, Answer: "17", Subject: "Math", Difficulty: "Hard"},
	// Math(extreme)
	{Text: "What is the value of 12^2 + 5^2?", Choices: []string{"169", "154", "149"}, Answer: "169", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the square root of 2025?", Choices: []string{"45", "40", "50"}, Answer: "45", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the result of (8 × 7) ÷ (2 + 2)?", Choices: []string{"14", "13", "15"}, Answer: "14", Subject: "Math", Difficulty: "Extreme"},
	{Text: "If x + y = 10 and x - y = 4, what is x?", Choices: []string{"7", "6", "5"}, Answer: "7", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the factorial of 5?", Choices: []string{"120", "60", "24"}, Answer: "120", Subject: "Math", Difficulty: "Extreme"},
	{Text: "Solve: (3^3 + 2^4) × 2", Choices: []string{"98", "100", "88"}, Answer: "98", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the derivative of 3x^2?", Choices: []string{"6x", "3x", "2x"}, Answer: "6x", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the area of a circle with radius 7?", Choices: []string{"154", "144", "132"}, Answer: "154", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is 111 × 111?", Choices: []string{"12321", "11111", "12221"}, Answer: "12321", Subject: "Math", Difficulty: "Extreme"},
	{Text: "Solve for x: 2x + 3 = 17", Choices: []string{"7", "6", "8"}, Answer: "7", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the sum of interior angles of a decagon?", Choices: []string{"1440", "1260", "1080"}, Answer: "1440", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is log₁₀(1000)?", Choices: []string{"3", "2", "1"}, Answer: "3", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the integral of x dx?", Choices: []string{"x^2 / 2 + C", "x^2 + C", "2x + C"}, Answer: "x^2 / 2 + C", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the 10th Fibonacci number?", Choices: []string{"55", "34", "89"}, Answer: "55", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is 2 to the power of 10?", Choices: []string{"1024", "1000", "512"}, Answer: "1024", Subject: "Math", Difficulty: "Extreme"},
	{Text: "How many primes are there between 1 and 20?", Choices: []string{"8", "7", "9"}, Answer: "8", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the inverse of 5/2?", Choices: []string{"2/5", "1/5", "5/1"}, Answer: "2/5", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the solution of x^2 - 4x + 4 = 0?", Choices: []string{"x = 2", "x = 4", "x = -2"}, Answer: "x = 2", Subject: "Math", Difficulty: "Extreme"},
	{Text: "Convert binary 1010 to decimal.", Choices: []string{"10", "12", "8"}, Answer: "10", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is sin(90°)?", Choices: []string{"1", "0", "0.5"}, Answer: "1", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the cube root of 729?", Choices: []string{"9", "8", "7"}, Answer: "9", Subject: "Math", Difficulty: "Extreme"},
	{Text: "Solve for x: x/3 = 7", Choices: []string{"21", "24", "18"}, Answer: "21", Subject: "Math", Difficulty: "Extreme"},
	{Text: "If a = 2 and b = 3, what is ab^2?", Choices: []string{"18", "12", "16"}, Answer: "18", Subject: "Math", Difficulty: "Extreme"},
	{Text: "What is the least common multiple of 6 and 8?", Choices: []string{"24", "48", "12"}, Answer: "24", Subject: "Math", Difficulty: "Extreme"},
	{Text: "If a triangle has sides 3, 4, and 5, what type is it?", Choices: []string{"Right", "Acute", "Obtuse"}, Answer: "Right", Subject: "Math", Difficulty: "Extreme"},
	{Text: "Evaluate: (4 + 5) × (6 - 2)", Choices: []string{"36", "32", "28"}, Answer: "36", Subject: "Math", Difficulty: "Extreme"},
	// English(easy)
	{Text: "What is the opposite of 'big'?", Choices: []string{"Small", "Tall", "Long"}, Answer: "Small", Subject: "English", Difficulty: "Easy"},
	{Text: "Which word is a noun?", Choices: []string{"Run", "Cat", "Blue"}, Answer: "Cat", Subject: "English", Difficulty: "Easy"},
	{Text: "Which one is a vowel?", Choices: []string{"B", "E", "T"}, Answer: "E", Subject: "English", Difficulty: "Easy"},
	{Text: "What sound does the letter 'B' make?", Choices: []string{"Buh", "Kuh", "Duh"}, Answer: "Buh", Subject: "English", Difficulty: "Easy"},
	{Text: "Which word rhymes with 'cat'?", Choices: []string{"Dog", "Hat", "Pig"}, Answer: "Hat", Subject: "English", Difficulty: "Easy"},
	{Text: "What is the correct article: '___ apple'?", Choices: []string{"A", "An", "The"}, Answer: "An", Subject: "English", Difficulty: "Easy"},
	{Text: "What comes first in the alphabet?", Choices: []string{"A", "C", "B"}, Answer: "A", Subject: "English", Difficulty: "Easy"},
	{Text: "Which is a color word?", Choices: []string{"Red", "Run", "Rat"}, Answer: "Red", Subject: "English", Difficulty: "Easy"},
	{Text: "What is the plural of 'dog'?", Choices: []string{"Dog", "Dogs", "Doges"}, Answer: "Dogs", Subject: "English", Difficulty: "Easy"},
	{Text: "Which is a question word?", Choices: []string{"Where", "Car", "Fast"}, Answer: "Where", Subject: "English", Difficulty: "Easy"},
	{Text: "Which sentence is correct?", Choices: []string{"He am happy.", "He is happy.", "He are happy."}, Answer: "He is happy.", Subject: "English", Difficulty: "Easy"},
	// English(Medium)
	{Text: "Which word is a noun?", Choices: []string{"run", "happy", "apple"}, Answer: "apple", Subject: "English", Difficulty: "Medium"},
	{Text: "What is the opposite of 'big'?", Choices: []string{"large", "small", "huge"}, Answer: "small", Subject: "English", Difficulty: "Medium"},
	{Text: "Which word begins with the letter 'B'?", Choices: []string{"cat", "bat", "apple"}, Answer: "bat", Subject: "English", Difficulty: "Medium"},
	{Text: "Choose the correct plural: one cat, two ___", Choices: []string{"cat", "cats", "cates"}, Answer: "cats", Subject: "English", Difficulty: "Medium"},
	{Text: "What do you do with your eyes?", Choices: []string{"hear", "see", "smell"}, Answer: "see", Subject: "English", Difficulty: "Medium"},
	{Text: "Which one is a color?", Choices: []string{"red", "run", "rat"}, Answer: "red", Subject: "English", Difficulty: "Medium"},
	{Text: "Which word is a verb?", Choices: []string{"sleep", "blue", "table"}, Answer: "sleep", Subject: "English", Difficulty: "Medium"},
	{Text: "Choose the correct word: The dog is ___ the box.", Choices: []string{"on", "under", "in"}, Answer: "in", Subject: "English", Difficulty: "Medium"},
	{Text: "What sound does 'ch' make in 'chicken'?", Choices: []string{"sh", "ch", "k"}, Answer: "ch", Subject: "English", Difficulty: "Medium"},
	{Text: "Which sentence is correct?", Choices: []string{"He run fast.", "He runs fast.", "He running fast."}, Answer: "He runs fast.", Subject: "English", Difficulty: "Medium"},
	{Text: "What do we call the name of a person?", Choices: []string{"adjective", "noun", "verb"}, Answer: "noun", Subject: "English", Difficulty: "Medium"},
	{Text: "Pick the correct word: I ___ a book.", Choices: []string{"am", "has", "have"}, Answer: "have", Subject: "English", Difficulty: "Medium"},
	{Text: "What is the past tense of 'jump'?", Choices: []string{"jumped", "jumping", "jumps"}, Answer: "jumped", Subject: "English", Difficulty: "Medium"},
	{Text: "Which one is a question word?", Choices: []string{"blue", "what", "book"}, Answer: "what", Subject: "English", Difficulty: "Medium"},
	{Text: "Which word rhymes with 'cake'?", Choices: []string{"make", "cat", "cup"}, Answer: "make", Subject: "English", Difficulty: "Medium"},
	{Text: "Choose the correct article: ___ apple is red.", Choices: []string{"A", "An", "The"}, Answer: "An", Subject: "English", Difficulty: "Medium"},
	// English(Hard)
	{Text: "Which word is a noun?", Choices: []string{"run", "happy", "cat"}, Answer: "cat", Subject: "English", Difficulty: "Hard"},
	{Text: "What is the opposite of 'cold'?", Choices: []string{"hot", "wet", "soft"}, Answer: "hot", Subject: "English", Difficulty: "Hard"},
	{Text: "Choose the correct sentence.", Choices: []string{"He go to school.", "He goes to school.", "He going to school."}, Answer: "He goes to school.", Subject: "English", Difficulty: "Hard"},
	{Text: "Which word rhymes with 'hat'?", Choices: []string{"pen", "rat", "dog"}, Answer: "rat", Subject: "English", Difficulty: "Hard"},
	{Text: "What is the past tense of 'jump'?", Choices: []string{"jumped", "jumping", "jumps"}, Answer: "jumped", Subject: "English", Difficulty: "Hard"},
	{Text: "Which of these is an adjective?", Choices: []string{"quick", "run", "boy"}, Answer: "quick", Subject: "English", Difficulty: "Hard"},
	{Text: "Choose the correct word: The bird is ___ the cage.", Choices: []string{"in", "at", "by"}, Answer: "in", Subject: "English", Difficulty: "Hard"},
	{Text: "Which sentence uses capital letters correctly?", Choices: []string{"i like ice cream.", "I Like Ice Cream.", "I like ice cream."}, Answer: "I like ice cream.", Subject: "English", Difficulty: "Hard"},
	{Text: "Which one is a question?", Choices: []string{"She is my sister.", "Is she your sister?", "She your sister."}, Answer: "Is she your sister?", Subject: "English", Difficulty: "Hard"},
	{Text: "Choose the word that starts with a vowel.", Choices: []string{"apple", "ball", "cat"}, Answer: "apple", Subject: "English", Difficulty: "Hard"},
	{Text: "Which word has the same beginning sound as 'sun'?", Choices: []string{"hat", "sand", "car"}, Answer: "sand", Subject: "English", Difficulty: "Hard"},
	{Text: "Which is a proper noun?", Choices: []string{"city", "man", "Zamboanga"}, Answer: "Zamboanga", Subject: "English", Difficulty: "Hard"},
	{Text: "What is the plural of 'baby'?", Choices: []string{"babys", "babies", "babes"}, Answer: "babies", Subject: "English", Difficulty: "Hard"},
	{Text: "Which word is a verb?", Choices: []string{"happy", "sleep", "blue"}, Answer: "sleep", Subject: "English", Difficulty: "Hard"},
	{Text: "Which word completes the sentence: She is ___ to the music.", Choices: []string{"listen", "listens", "listening"}, Answer: "listening", Subject: "English", Difficulty: "Hard"},
	{Text: "Choose the correct punctuation for this sentence: What is your name", Choices: []string{"?", ".", "!"}, Answer: "?", Subject: "English", Difficulty: "Hard"},
	{Text: "Which word means the same as 'big'?", Choices: []string{"small", "huge", "thin"}, Answer: "huge", Subject: "English", Difficulty: "Hard"},
	{Text: "What is the correct article: ___ elephant is big.", Choices: []string{"A", "An", "The"}, Answer: "An", Subject: "English", Difficulty: "Hard"},
	{Text: "Choose the word with a silent letter.", Choices: []string{"knee", "sun", "dog"}, Answer: "knee", Subject: "English", Difficulty: "Hard"},
	{Text: "Which one is a compound word?", Choices: []string{"sunshine", "sun", "shine"}, Answer: "sunshine", Subject: "English", Difficulty: "Hard"},
	{Text: "Choose the correct homophone: I went ___ the store.", Choices: []string{"to", "two", "too"}, Answer: "to", Subject: "English", Difficulty: "Hard"},
	// English(Extreme)
	{Text: "Which sentence uses the correct past tense?", Choices: []string{"She go to school.", "She went to school.", "She going to school."}, Answer: "She went to school.", Subject: "English", Difficulty: "Extreme"},
	{Text: "What is the opposite of 'begin'?", Choices: []string{"End", "Start", "Continue"}, Answer: "End", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which word is a noun?", Choices: []string{"Run", "Blue", "Chair"}, Answer: "Chair", Subject: "English", Difficulty: "Extreme"},
	{Text: "Choose the correct plural: 'mouse'", Choices: []string{"Mouses", "Mice", "Mouse"}, Answer: "Mice", Subject: "English", Difficulty: "Extreme"},
	{Text: "What is a synonym for 'happy'?", Choices: []string{"Sad", "Joyful", "Angry"}, Answer: "Joyful", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which sentence is a question?", Choices: []string{"Where are you going", "Where are you going?", "Where you going."}, Answer: "Where are you going?", Subject: "English", Difficulty: "Extreme"},
	{Text: "What does 'predict' mean?", Choices: []string{"To look back", "To say what will happen", "To fix something"}, Answer: "To say what will happen", Subject: "English", Difficulty: "Extreme"},
	{Text: "Choose the correct contraction: 'She is'", Choices: []string{"She's", "Shes", "She is'"}, Answer: "She's", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which is an adjective?", Choices: []string{"Quickly", "Beautiful", "Run"}, Answer: "Beautiful", Subject: "English", Difficulty: "Extreme"},
	{Text: "What type of sentence is 'Wow! That's amazing!'", Choices: []string{"Declarative", "Interrogative", "Exclamatory"}, Answer: "Exclamatory", Subject: "English", Difficulty: "Extreme"},
	{Text: "What punctuation ends a question?", Choices: []string{"!", ".", "?"}, Answer: "?", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which word means the same as 'tiny'?", Choices: []string{"Large", "Small", "Wide"}, Answer: "Small", Subject: "English", Difficulty: "Extreme"},
	{Text: "Choose the correct possessive form: The toys of the dog", Choices: []string{"Dogs toy", "Dog's toys", "Dogs' toys"}, Answer: "Dog's toys", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which word is a verb?", Choices: []string{"Table", "Jump", "Blue"}, Answer: "Jump", Subject: "English", Difficulty: "Extreme"},
	{Text: "Choose the correct sentence.", Choices: []string{"He goed to school.", "He went to school.", "He going to school."}, Answer: "He went to school.", Subject: "English", Difficulty: "Extreme"},
	{Text: "What does 'antonym' mean?", Choices: []string{"A word that means the same", "A word that means the opposite", "A describing word"}, Answer: "A word that means the opposite", Subject: "English", Difficulty: "Extreme"},
	{Text: "Pick the correct homophone: 'Their going to the park.'", Choices: []string{"Their", "There", "They're"}, Answer: "They're", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which word is spelled correctly?", Choices: []string{"Becaus", "Because", "Becuz"}, Answer: "Because", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which word best completes the sentence: 'She ___ to the store.'", Choices: []string{"go", "went", "going"}, Answer: "went", Subject: "English", Difficulty: "Extreme"},
	{Text: "What part of speech is the word 'quickly'?", Choices: []string{"Adjective", "Noun", "Adverb"}, Answer: "Adverb", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which is a compound word?", Choices: []string{"Sunlight", "Light", "Sun"}, Answer: "Sunlight", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which of these is a proper noun?", Choices: []string{"city", "store", "London"}, Answer: "London", Subject: "English", Difficulty: "Extreme"},
	{Text: "What do quotation marks show?", Choices: []string{"An action", "A question", "Someone is speaking"}, Answer: "Someone is speaking", Subject: "English", Difficulty: "Extreme"},
	{Text: "Choose the sentence with correct punctuation.", Choices: []string{"what time is it", "What time is it.", "What time is it?"}, Answer: "What time is it?", Subject: "English", Difficulty: "Extreme"},
	{Text: "What does 'prefix' mean?", Choices: []string{"A word at the end", "A word at the start", "A word in the middle"}, Answer: "A word at the start", Subject: "English", Difficulty: "Extreme"},
	{Text: "Which sentence uses 'there' correctly?", Choices: []string{"There going to the mall.", "The dog is over there.", "There house is big."}, Answer: "The dog is over there.", Subject: "English", Difficulty: "Extreme"},
	// Science(Easy)
	{Text: "What do we breathe in to live?", Choices: []string{"Water", "Oxygen", "Smoke"}, Answer: "Oxygen", Subject: "Science", Difficulty: "Easy"},
	{Text: "What do plants need to grow?", Choices: []string{"Milk", "Sunlight", "Juice"}, Answer: "Sunlight", Subject: "Science", Difficulty: "Easy"},
	{Text: "What is the color of the sky on a clear day?", Choices: []string{"Blue", "Green", "Red"}, Answer: "Blue", Subject: "Science", Difficulty: "Easy"},
	{Text: "Which of these is a sense organ?", Choices: []string{"Heart", "Ear", "Liver"}, Answer: "Ear", Subject: "Science", Difficulty: "Easy"},
	{Text: "Which part of the plant is green and makes food?", Choices: []string{"Stem", "Leaf", "Root"}, Answer: "Leaf", Subject: "Science", Difficulty: "Easy"},
	{Text: "What do fish use to breathe?", Choices: []string{"Nose", "Gills", "Mouth"}, Answer: "Gills", Subject: "Science", Difficulty: "Easy"},
	{Text: "What do bees make?", Choices: []string{"Milk", "Honey", "Bread"}, Answer: "Honey", Subject: "Science", Difficulty: "Easy"},
	{Text: "Which animal can fly?", Choices: []string{"Cat", "Dog", "Bird"}, Answer: "Bird", Subject: "Science", Difficulty: "Easy"},
	{Text: "What do we use to see things?", Choices: []string{"Nose", "Eyes", "Ears"}, Answer: "Eyes", Subject: "Science", Difficulty: "Easy"},
	{Text: "What do you drink when you are thirsty?", Choices: []string{"Soda", "Juice", "Water"}, Answer: "Water", Subject: "Science", Difficulty: "Easy"},
	{Text: "What is the sun?", Choices: []string{"A planet", "A star", "A moon"}, Answer: "A star", Subject: "Science", Difficulty: "Easy"},
	// Science(Medium)
	{Text: "What do plants need to grow?", Choices: []string{"Milk", "Sunlight", "Sugar"}, Answer: "Sunlight", Subject: "Science", Difficulty: "Medium"},
	{Text: "Which part of the body helps us see?", Choices: []string{"Ears", "Eyes", "Nose"}, Answer: "Eyes", Subject: "Science", Difficulty: "Medium"},
	{Text: "What do we breathe in to stay alive?", Choices: []string{"Water", "Oxygen", "Smoke"}, Answer: "Oxygen", Subject: "Science", Difficulty: "Medium"},
	{Text: "What is the color of healthy leaves?", Choices: []string{"Brown", "Yellow", "Green"}, Answer: "Green", Subject: "Science", Difficulty: "Medium"},
	{Text: "Which animal lays eggs?", Choices: []string{"Dog", "Cat", "Chicken"}, Answer: "Chicken", Subject: "Science", Difficulty: "Medium"},
	{Text: "What do fish use to swim?", Choices: []string{"Legs", "Fins", "Wings"}, Answer: "Fins", Subject: "Science", Difficulty: "Medium"},
	{Text: "Where does the sun go at night?", Choices: []string{"It sleeps", "It hides", "It moves to the other side of the Earth"}, Answer: "It moves to the other side of the Earth", Subject: "Science", Difficulty: "Medium"},
	{Text: "What do we use to smell things?", Choices: []string{"Mouth", "Hands", "Nose"}, Answer: "Nose", Subject: "Science", Difficulty: "Medium"},
	{Text: "What helps us hear sounds?", Choices: []string{"Eyes", "Ears", "Hands"}, Answer: "Ears", Subject: "Science", Difficulty: "Medium"},
	{Text: "Which of these is a living thing?", Choices: []string{"Rock", "Tree", "Car"}, Answer: "Tree", Subject: "Science", Difficulty: "Medium"},
	{Text: "What do roots do for a plant?", Choices: []string{"Help it breathe", "Hold it in the soil", "Make flowers"}, Answer: "Hold it in the soil", Subject: "Science", Difficulty: "Medium"},
	{Text: "Which of these can fly?", Choices: []string{"Dog", "Bird", "Snake"}, Answer: "Bird", Subject: "Science", Difficulty: "Medium"},
	{Text: "What happens when ice is left in the sun?", Choices: []string{"It gets bigger", "It melts", "It turns into dust"}, Answer: "It melts", Subject: "Science", Difficulty: "Medium"},
	{Text: "Which sense helps you feel a soft teddy bear?", Choices: []string{"Sight", "Touch", "Taste"}, Answer: "Touch", Subject: "Science", Difficulty: "Medium"},
	{Text: "What covers and protects your body?", Choices: []string{"Bones", "Skin", "Hair"}, Answer: "Skin", Subject: "Science", Difficulty: "Medium"},
	{Text: "Which of these grows from a seed?", Choices: []string{"Table", "Flower", "Toy"}, Answer: "Flower", Subject: "Science", Difficulty: "Medium"},
	// Science(Hard)
	{Text: "Which part of the plant makes food?", Choices: []string{"Roots", "Leaves", "Stem"}, Answer: "Leaves", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do humans need to breathe?", Choices: []string{"Oxygen", "Carbon Dioxide", "Water"}, Answer: "Oxygen", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which of these is not a living thing?", Choices: []string{"Tree", "Rock", "Dog"}, Answer: "Rock", Subject: "Science", Difficulty: "Hard"},
	{Text: "Where does a fish live?", Choices: []string{"Air", "Water", "Land"}, Answer: "Water", Subject: "Science", Difficulty: "Hard"},
	{Text: "What helps humans see?", Choices: []string{"Ears", "Eyes", "Hands"}, Answer: "Eyes", Subject: "Science", Difficulty: "Hard"},
	{Text: "What does the Sun give us?", Choices: []string{"Light", "Food", "Water"}, Answer: "Light", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which animal lays eggs?", Choices: []string{"Cat", "Chicken", "Dog"}, Answer: "Chicken", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which part of the body helps us smell?", Choices: []string{"Ears", "Nose", "Mouth"}, Answer: "Nose", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do plants need to grow?", Choices: []string{"Sunlight", "Sugar", "Wind"}, Answer: "Sunlight", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do we call water that falls from the sky?", Choices: []string{"Snow", "Rain", "Sun"}, Answer: "Rain", Subject: "Science", Difficulty: "Hard"},
	{Text: "What is the largest organ of the body?", Choices: []string{"Heart", "Skin", "Liver"}, Answer: "Skin", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which object is in the sky during the day?", Choices: []string{"Moon", "Stars", "Sun"}, Answer: "Sun", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do bees make?", Choices: []string{"Milk", "Honey", "Sugar"}, Answer: "Honey", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which of these is a gas?", Choices: []string{"Ice", "Water", "Air"}, Answer: "Air", Subject: "Science", Difficulty: "Hard"},
	{Text: "What happens to water when it is frozen?", Choices: []string{"It becomes ice", "It disappears", "It boils"}, Answer: "It becomes ice", Subject: "Science", Difficulty: "Hard"},
	{Text: "What part of the tree is under the ground?", Choices: []string{"Trunk", "Leaves", "Roots"}, Answer: "Roots", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which of the following can grow?", Choices: []string{"Rock", "Flower", "Spoon"}, Answer: "Flower", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do you call a young cat?", Choices: []string{"Puppy", "Kitten", "Cub"}, Answer: "Kitten", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do we call animals that live in water?", Choices: []string{"Insects", "Fish", "Birds"}, Answer: "Fish", Subject: "Science", Difficulty: "Hard"},
	{Text: "Which one of these animals can fly?", Choices: []string{"Bat", "Dog", "Frog"}, Answer: "Bat", Subject: "Science", Difficulty: "Hard"},
	{Text: "What do plants give off that helps us breathe?", Choices: []string{"Oxygen", "Smoke", "Dust"}, Answer: "Oxygen", Subject: "Science", Difficulty: "Hard"},
	// Science(Extreme)
	{Text: "What part of the plant makes food?", Choices: []string{"Leaf", "Root", "Stem"}, Answer: "Leaf", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What do humans need to breathe?", Choices: []string{"Oxygen", "Carbon dioxide", "Water"}, Answer: "Oxygen", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What is the hardest part of your body?", Choices: []string{"Skin", "Bone", "Tooth"}, Answer: "Tooth", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which planet is closest to the Sun?", Choices: []string{"Earth", "Mars", "Mercury"}, Answer: "Mercury", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What gas do plants give off?", Choices: []string{"Oxygen", "Nitrogen", "Hydrogen"}, Answer: "Oxygen", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What do you call water when it turns into gas?", Choices: []string{"Ice", "Steam", "Snow"}, Answer: "Steam", Subject: "Science", Difficulty: "Extreme"},
	{Text: "How many legs does an insect have?", Choices: []string{"4", "6", "8"}, Answer: "6", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What is the function of roots in plants?", Choices: []string{"Make food", "Absorb water", "Help breathe"}, Answer: "Absorb water", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which sense organ helps us smell?", Choices: []string{"Eye", "Ear", "Nose"}, Answer: "Nose", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What is the smallest unit of life?", Choices: []string{"Organ", "Cell", "Tissue"}, Answer: "Cell", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which part of the body helps pump blood?", Choices: []string{"Lung", "Heart", "Stomach"}, Answer: "Heart", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What helps plants grow toward light?", Choices: []string{"Phototropism", "Photosynthesis", "Evaporation"}, Answer: "Phototropism", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which liquid helps digest food in the stomach?", Choices: []string{"Saliva", "Acid", "Water"}, Answer: "Acid", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What happens when the sun sets?", Choices: []string{"Morning", "Evening", "Afternoon"}, Answer: "Evening", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What tool do you use to look at stars?", Choices: []string{"Microscope", "Binoculars", "Telescope"}, Answer: "Telescope", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which planet has rings?", Choices: []string{"Venus", "Saturn", "Mars"}, Answer: "Saturn", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What is rain made of?", Choices: []string{"Dust", "Water", "Gas"}, Answer: "Water", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What makes the moon shine at night?", Choices: []string{"It glows", "It reflects sunlight", "It burns gas"}, Answer: "It reflects sunlight", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What animal lays eggs and can fly?", Choices: []string{"Bat", "Eagle", "Dog"}, Answer: "Eagle", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Where do fish get oxygen?", Choices: []string{"Air", "Water", "Sand"}, Answer: "Water", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What state of matter is steam?", Choices: []string{"Solid", "Liquid", "Gas"}, Answer: "Gas", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which part of the plant holds it upright?", Choices: []string{"Leaf", "Stem", "Flower"}, Answer: "Stem", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Which animal has scales and lays eggs?", Choices: []string{"Dog", "Lizard", "Frog"}, Answer: "Lizard", Subject: "Science", Difficulty: "Extreme"},
	{Text: "Why do we see lightning before thunder?", Choices: []string{"Light is faster than sound", "Sound is louder", "Thunder comes first"}, Answer: "Light is faster than sound", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What are clouds made of?", Choices: []string{"Dust", "Water droplets", "Air"}, Answer: "Water droplets", Subject: "Science", Difficulty: "Extreme"},
	{Text: "What do birds use to fly?", Choices: []string{"Feet", "Wings", "Beak"}, Answer: "Wings", Subject: "Science", Difficulty: "Extreme"},
	// Filipino(Easy)
	{Text: "Ano ang pambansang prutas ng Pilipinas?", Choices: []string{"Mangga", "Saging", "Pakwan"}, Answer: "Mangga", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Ano ang kulay ng dahon?", Choices: []string{"Pula", "Berde", "Dilaw"}, Answer: "Berde", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Anong tunog ang ginagawa ng aso?", Choices: []string{"Moo", "Kokak", "Aw aw"}, Answer: "Aw aw", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Anong bahagi ng katawan ang ginagamit sa paglakad?", Choices: []string{"Kamay", "Paa", "Ulo"}, Answer: "Paa", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Ano ang iniinom ng sanggol?", Choices: []string{"Tubig", "Kape", "Gatas"}, Answer: "Gatas", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Anong hayop ang may mahabang leeg?", Choices: []string{"Aso", "Giraffe", "Pusa"}, Answer: "Giraffe", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Ano ang tawag sa tatay ng iyong nanay?", Choices: []string{"Lolo", "Tito", "Kuya"}, Answer: "Lolo", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Anong araw ang kasunod ng Lunes?", Choices: []string{"Sabado", "Martes", "Linggo"}, Answer: "Martes", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Ano ang tawag sa kulay ng langit tuwing umaga?", Choices: []string{"Itim", "Berde", "Bughaw"}, Answer: "Bughaw", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Saan pumapasok ang mga bata para mag-aral?", Choices: []string{"Palengke", "Sinehan", "Paaralan"}, Answer: "Paaralan", Subject: "Filipino", Difficulty: "Easy"},
	{Text: "Anong prutas ang may maraming mata?", Choices: []string{"Saging", "Pinya", "Mangga"}, Answer: "Pinya", Subject: "Filipino", Difficulty: "Easy"},
	// Filipino(Medium)
	{Text: "Ano ang unang titik ng salitang 'aso'?", Choices: []string{"a", "o", "s"}, Answer: "a", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Ano ang tawag sa larawan ng araw, ulap, at ulan?", Choices: []string{"Panahon", "Pagkain", "Hayop"}, Answer: "Panahon", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Ano ang kabaligtaran ng salitang 'malaki'?", Choices: []string{"mahaba", "maliit", "mabigat"}, Answer: "maliit", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Piliin ang salitang nagsisimula sa letrang 'b'.", Choices: []string{"aso", "bola", "gatas"}, Answer: "bola", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Anong tunog ang naririnig sa dulo ng salitang 'gabi'?", Choices: []string{"a", "i", "g"}, Answer: "i", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Alin sa mga sumusunod ang bahagi ng katawan?", Choices: []string{"kamay", "mesa", "bintana"}, Answer: "kamay", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Pumili ng tamang sagot: Ang langit ay ____.", Choices: []string{"berde", "asul", "dilaw"}, Answer: "asul", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Ano ang kasunod ng titik 'k' sa alpabeto?", Choices: []string{"j", "l", "m"}, Answer: "l", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Alin sa mga ito ang hayop?", Choices: []string{"kabayo", "sapatos", "upuan"}, Answer: "kabayo", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Ano ang tamang pantig ng salitang 'bahay'?", Choices: []string{"ba-hay", "bah-ay", "baha-y"}, Answer: "ba-hay", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Piliin ang salitang pareho ang tunog sa 'bata'.", Choices: []string{"mata", "mesa", "pusa"}, Answer: "mata", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Ano ang dapat gamitin sa dulo ng pangungusap?", Choices: []string{"tuldok", "kwit", "kudlit"}, Answer: "tuldok", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Saan ginagamit ang walis?", Choices: []string{"sa pagkain", "sa pagsusulat", "sa paglinis"}, Answer: "sa paglinis", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Alin sa mga ito ang hindi kulay?", Choices: []string{"pula", "dahon", "asul"}, Answer: "dahon", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Pumili ng salitang tumutukoy sa prutas.", Choices: []string{"saging", "lapis", "silya"}, Answer: "saging", Subject: "Filipino", Difficulty: "Medium"},
	{Text: "Ano ang salitang angkop sa 'Ang ibon ay ___ sa langit.'?", Choices: []string{"lumilipad", "kumakain", "natutulog"}, Answer: "lumilipad", Subject: "Filipino", Difficulty: "Medium"},
	// Filipino(Hard)
	{Text: "Ano ang tamang baybay ng salitang 'maganda'?", Choices: []string{"maganda", "magannda", "magnda"}, Answer: "maganda", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang kasalungat ng salitang 'malaki'?", Choices: []string{"mataas", "maliit", "mahaba"}, Answer: "maliit", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Pumili ng pangngalan: Si Ana ay kumain ng mangga.", Choices: []string{"kumain", "Ana", "mangga"}, Answer: "Ana", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Anong bahagi ng katawan ang ginagamit sa pagdinig?", Choices: []string{"mata", "tainga", "ilong"}, Answer: "tainga", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Anong tunog ang nagsisimula sa titik 'B'?", Choices: []string{"aso", "bola", "gatas"}, Answer: "bola", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Alin ang tamang gamit ng 'ng'?", Choices: []string{"Kumain ng saging.", "Ng umaga ay malamig.", "Ng bahay ay malaki."}, Answer: "Kumain ng saging.", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Piliin ang salitang may diptonggo.", Choices: []string{"bata", "gabi", "araw"}, Answer: "araw", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang sagot sa bugtong: Isang balong malalim, punong-puno ng patalim?", Choices: []string{"bibig", "balon", "ilong"}, Answer: "bibig", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Alin ang salitang pantig na 'ka-la-ba-sa'?", Choices: []string{"kalabasa", "kamatis", "kahel"}, Answer: "kalabasa", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Piliin ang tamang gamit ng 'ang'.", Choices: []string{"Ang bahay ay malaki.", "Bahay ang malaki.", "Malaki ang bahay."}, Answer: "Ang bahay ay malaki.", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Alin sa mga sumusunod ang pang-uri?", Choices: []string{"kumain", "mataas", "siya"}, Answer: "mataas", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Piliin ang tamang baybay: anák, anak, ana'k", Choices: []string{"anak", "anák", "ana'k"}, Answer: "anak", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang kasalungat ng 'maingay'?", Choices: []string{"matahimik", "masaya", "malakas"}, Answer: "matahimik", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Anong bahagi ng pangungusap ang 'umalis si kuya'?", Choices: []string{"simuno", "panaguri", "pangatnig"}, Answer: "panaguri", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Alin ang panghalip panao?", Choices: []string{"siya", "bata", "bahay"}, Answer: "siya", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang kahulugan ng 'matimtiman'?", Choices: []string{"malakas", "matahimik", "masunurin"}, Answer: "masunurin", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Piliin ang tambalang salita.", Choices: []string{"bahaghari", "gabi", "gatas"}, Answer: "bahaghari", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang tawag sa tunog ng aso?", Choices: []string{"meow", "tiktilaok", "aw-aw"}, Answer: "aw-aw", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang dapat gamitin sa katapusan ng tanong?", Choices: []string{"tuldok", "tandang padamdam", "tandang pananong"}, Answer: "tandang pananong", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Piliin ang tama: Si Jose ay ______ ng kendi.", Choices: []string{"kumain", "kumakain", "kinain"}, Answer: "kumain", Subject: "Filipino", Difficulty: "Hard"},
	{Text: "Ano ang sagot sa bugtong: May puno walang bunga, may dahon walang sanga?", Choices: []string{"payong", "mesa", "libro"}, Answer: "payong", Subject: "Filipino", Difficulty: "Hard"},
	// Filipino(Extreme)
	{Text: "Ano ang kabaligtaran ng 'masaya'?", Choices: []string{"Malungkot", "Masigla", "Masarap"}, Answer: "Malungkot", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Anong tunog ang unang maririnig sa salitang 'kabayo'?", Choices: []string{"Ka", "Ba", "Yo"}, Answer: "Ka", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang tawag sa larawang gumagamit ng salita?", Choices: []string{"Tula", "Pabula", "Kuwento"}, Answer: "Kuwento", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Piliin ang tamang gamit ng 'ng': Ang bata ___ umiyak ay nawalan ng laruan.", Choices: []string{"na", "ng", "nang"}, Answer: "ng", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Alin ang salitang kilos?", Choices: []string{"Takbo", "Mesa", "Gabi"}, Answer: "Takbo", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang kasunod ng titik E sa alpabeto?", Choices: []string{"F", "D", "G"}, Answer: "F", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Piliin ang wastong baybay: 'Mga taong nag-aaral'", Choices: []string{"Estudyante", "Estudiyante", "Estudyanti"}, Answer: "Estudyante", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang ibig sabihin ng 'matimtiman'?", Choices: []string{"Maingay", "Tahimik", "Maayos"}, Answer: "Tahimik", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Alin sa mga ito ang pantig ng salitang 'umaga'?", Choices: []string{"u-ma-ga", "um-a-ga", "uma-ga"}, Answer: "u-ma-ga", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ang salitang 'maganda' ay kabaligtaran ng?", Choices: []string{"Pangit", "Maayos", "Mabait"}, Answer: "Pangit", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Anong uri ng hayop si 'Pagong' sa kuwentong 'Pagong at Matsing'?", Choices: []string{"Reptilya", "Isda", "Amphibian"}, Answer: "Reptilya", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang ibig sabihin ng salitang 'masigasig'?", Choices: []string{"Tamad", "Masipag", "Malakas"}, Answer: "Masipag", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Alin ang hindi kabilang: Bola, Aklat, Isda, Lapís?", Choices: []string{"Isda", "Bola", "Lapis"}, Answer: "Isda", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Saan ginagamit ang salitang 'po' at 'opo'?", Choices: []string{"Sa kaibigan", "Sa bata", "Sa nakatatanda"}, Answer: "Sa nakatatanda", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ang 'Aso ay tumatahol.' Ano ang pandiwa?", Choices: []string{"Aso", "Tumatahol", "Ay"}, Answer: "Tumatahol", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ang salitang 'bituin' ay tumutukoy sa?", Choices: []string{"Hayop", "Bagay", "Kalangitan"}, Answer: "Kalangitan", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Anong bahagi ng katawan ang ginagamit sa pandinig?", Choices: []string{"Ilong", "Tenga", "Mata"}, Answer: "Tenga", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang kasingkahulugan ng 'mabilis'?", Choices: []string{"Mabagal", "Matulin", "Malakas"}, Answer: "Matulin", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang ibig sabihin ng 'panaginip'?", Choices: []string{"Gabi", "Isip", "Nasa isip habang natutulog"}, Answer: "Nasa isip habang natutulog", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Alin sa mga ito ang isang bugtong?", Choices: []string{"Maliit na bahay, puno ng halakhak", "May aso sa labas", "Nagmamadaling bata"}, Answer: "Maliit na bahay, puno ng halakhak", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang kasingkahulugan ng 'masaya'?", Choices: []string{"Malungkot", "Masigla", "Matapang"}, Answer: "Masigla", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Alin ang tamang gamit ng 'nang': Kumain siya ___ tahimik.", Choices: []string{"ng", "nang", "na"}, Answer: "nang", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang ibig sabihin ng 'umaga'?", Choices: []string{"Gabi", "Gitna ng araw", "Simula ng araw"}, Answer: "Simula ng araw", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Alin ang wastong sagot: Ako ay may alaga, ___ ay pusa.", Choices: []string{"Siya", "Ito", "Ito'y"}, Answer: "Ito", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang gamit ng bantas na tandang pananong ( ? )?", Choices: []string{"Sa tanong", "Sa utos", "Sa kuwento"}, Answer: "Sa tanong", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Ano ang salitang inuulit sa 'araw-araw'?", Choices: []string{"Araw", "Raw", "Wala"}, Answer: "Araw", Subject: "Filipino", Difficulty: "Extreme"},
	{Text: "Anong uri ng pangungusap ito: 'Pakibuksan ang pinto.'?", Choices: []string{"Pautos", "Patanong", "Pasalaysay"}, Answer: "Pautos", Subject: "Filipino", Difficulty: "Extreme"},
}

// NewQuiz creates a new quiz with sample questions
func NewQuiz() *Quiz {
	return &Quiz{Questions: sampleQuestions}
}

// ListSubjects returns all unique subjects
func (q *Quiz) ListSubjects() []string {
	subjectMap := make(map[string]struct{})
	for _, ques := range q.Questions {
		subjectMap[ques.Subject] = struct{}{}
	}
	subjects := make([]string, 0, len(subjectMap))
	for s := range subjectMap {
		subjects = append(subjects, s)
	}
	return subjects
}

// GetRandomQuestion returns a random question for a subject and difficulty
func (q *Quiz) GetRandomQuestion(subject, difficulty string) *Question {
	var filtered []Question
	for _, ques := range q.Questions {
		if strings.EqualFold(ques.Subject, subject) && strings.EqualFold(ques.Difficulty, difficulty) {
			filtered = append(filtered, ques)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(filtered))
	return &filtered[idx]
}

// CheckAnswer checks if the answer is correct (case-insensitive)
func (q *Quiz) CheckAnswer(ques *Question, answer string) bool {
	return strings.EqualFold(strings.TrimSpace(ques.Answer), strings.TrimSpace(answer))
}

func InsertUser(name string) (int64, error) {
	res, err := DB.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func InsertLeaderboard(userID int64, score, quests, boosts int, accuracy, bonus float64) error {
	_, err := DB.Exec(
		"INSERT INTO leaderboard (user_id, score, quests_completed, weapon_boosts, accuracy, bonus_success) VALUES (?, ?, ?, ?, ?, ?)",
		userID, score, quests, boosts, accuracy, bonus,
	)
	return err
}

func GetTopLeaderboard(limit int) ([]LeaderboardEntry, error) {
	rows, err := DB.Query(
		`SELECT u.name, l.score, l.quests_completed, l.weapon_boosts, l.accuracy, l.bonus_success
		 FROM leaderboard l
		 JOIN users u ON l.user_id = u.id
		 ORDER BY l.score DESC
		 LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		if err := rows.Scan(&e.PlayerName, &e.Score, &e.QuestsCompleted, &e.WeaponBoosts, &e.Accuracy, &e.BonusSuccess); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// Check if user has completed a difficulty for a subject
func HasCompletedDifficulty(userID int64, subject, difficulty string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND subject = ? AND difficulty = ?", userID, subject, difficulty).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Mark a difficulty as completed for a user and subject
func MarkDifficultyCompleted(userID int64, subject, difficulty string) error {
	_, err := DB.Exec("INSERT INTO user_progress (user_id, subject, difficulty) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE completed_at = CURRENT_TIMESTAMP", userID, subject, difficulty)
	return err
}

// GetTotalScore returns the total score from all subjects for a user
func GetTotalScore(userID int64) (int, error) {
	var totalScore int
	err := DB.QueryRow("SELECT COALESCE(SUM(score), 0) FROM leaderboard WHERE user_id = ?", userID).Scan(&totalScore)
	return totalScore, err
}
