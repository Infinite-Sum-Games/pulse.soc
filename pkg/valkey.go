package pkg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"github.com/redis/go-redis/v9"
)

var Valkey *redis.Client

const (
	Leaderboard = "leader-board-sset"
	CppRank     = "cpp-ranking-sset"
	JavaRank    = "java-ranking-sset"
	PyRank      = "python-ranking-sset"
	JsRank      = "javascript-ranking-sset"
	GoRank      = "go-ranking-sset"
	RustRank    = "rust-ranking-sset"
	ZigRank     = "zig-ranking-sset"
	FlutterRank = "flutter-ranking-sset"
	KotlinRank  = "kotlin-ranking-sset"
	HaskellRank = "haskell-ranking-sset"
)

func GetParticipantRank(username string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// RevRank returns position in descending order (0-based). If there is no
	// rank available then return -1. If there is an error then also return -1.
	// Incase of unavailability the error is redis.Nil which is bubbled up
	rank, err := Valkey.ZRevRank(ctx, Leaderboard, username).Result()
	if err != nil {
		return -1, err
	}
	return rank + 1, nil
}

// The ParticipantGlobal struct captures username and bounty for the leaderboard
// from the individual sorted sets for each language and returns only the top 2
type ParticipantGlobal struct {
	Username string `json:"github_username"`
	Bounty   string `json:"bounty"`
}

func GetLeaderboard() ([]ParticipantGlobal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results := []ParticipantGlobal{}
	participants, err := Valkey.ZRevRangeByScoreWithScores(ctx, Leaderboard,
		&redis.ZRangeBy{
			Min: "1",
			Max: "+inf",
		}).Result()
	if err != nil {
		return results, err
	}

	for _, participant := range participants {
		results = append(results, ParticipantGlobal{
			Username: fmt.Sprintf("%v", participant.Member),
			Bounty:   fmt.Sprintf("%d", int(participant.Score)),
		})
	}
	return results, nil
}

// The Participant struct captures participant username and pull_request_count
// from the individual sorted sets for each language and returns only the top 2
type Participant struct {
	Username string `json:"github_username"`
	Score    string `json:"pull_request_count"`
}

func GetTopParticipants() (map[string]map[string]Participant, error) {
	ssets := []string{
		CppRank, JavaRank,
		PyRank, JsRank,
		GoRank, RustRank,
		ZigRank, FlutterRank,
		KotlinRank, HaskellRank,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	results := make(map[string]map[string]Participant)
	var err error
	for _, sset := range ssets {
		name := strings.Split(sset, "-")[0] // extracting the language name
		members, err := Valkey.ZRevRangeByScoreWithScores(ctx, sset,
			&redis.ZRangeBy{
				Min:    "1",
				Max:    "+inf",
				Offset: 0,
				Count:  2,
			}).Result()

		// Setting up defaults
		results[name]["first_place"] = Participant{
			Username: "-",
			Score:    "-1",
		}
		results[name]["second_place"] = Participant{
			Username: "-",
			Score:    "-1",
		}

		// If unable to find the sorted-set due to Redis errors
		if err != nil {
			cmd.Log.Error("Error while trying to query sorted-set:"+sset, err)
			return nil, err
		}
		// There can be an edge case where only one participant has made PRs
		// in a particular language. Need to handle that case where there isn't
		// two entries
		if len(members) > 0 {
			results[name]["first_place"] = Participant{
				Username: fmt.Sprintf("%v", members[0].Member),
				Score:    fmt.Sprintf("%d", int(members[0].Score)),
			}
		}
		if len(members) > 1 {
			results[name]["second_place"] = Participant{
				Username: fmt.Sprintf("%v", members[1].Member),
				Score:    fmt.Sprintf("%d", int(members[1].Score)),
			}
		}
	}
	return results, err
}
