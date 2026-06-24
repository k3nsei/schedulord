package github

import "context"

type githubTokenKey struct{}

func WithGithubToken(parent context.Context, token string) context.Context {
	return context.WithValue(parent, githubTokenKey{}, token)
}

func GetGithubToken(ctx context.Context) (string, bool) {
	if v := ctx.Value(githubTokenKey{}); v != nil {
		if t, ok := v.(string); ok {
			return t, true
		}
	}
	return "", false
}
