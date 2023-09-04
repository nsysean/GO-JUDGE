package enums

import models "app/src/models"

var (
	VERDICT_ACCEPTED models.Verdict = models.Verdict{
		Prefix: "ok",
		Verdict: "Accepted",
	}

	VERDICT_PARTIAL models.Verdict = models.Verdict{
		Prefix: "points",
		Verdict: "Partial Score",
	}

	VERDICT_WA models.Verdict = models.Verdict{
		Prefix: "wrong answer",
		Verdict: "Wrong Answer",
	}

	VERDICT_WOF models.Verdict = models.Verdict{
		Prefix: "wrong output format",
		Verdict: "Presentation Error",
	}

	VERDICT_UEOF models.Verdict = models.Verdict{
		Prefix: "unexpected eof",
		Verdict: "Presentation Error",
	}

	VERDICT_TLE models.Verdict = models.Verdict{
		Verdict: "Time Limit Exceeded",
	}

	VERDICT_MLE models.Verdict = models.Verdict{
		Verdict: "Memory Limit Exceeded",
	}

	VERDICT_RE models.Verdict = models.Verdict{
		Verdict: "Runtime Error",
	}
)
