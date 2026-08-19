package config

import (
	"os"
	"strings"
	"testing"
)

type test struct {
	conf string
	user string
	want bool
}

var tests = []test{
	{"deny =all\n	allow =all\n", "user", false},
	{"deny = all\n	allow =user1,user\n", "user", true},
	{"deny =all\n	allow =user1,user\n", "user1", true},
	{"deny =all\n	allow =user1\n", "user", false},
	{"deny =user\n	allow =all\n", "user", false},
	{"deny =user1,user2\n	allow =all\n", "user", true},
	{"deny =user1,user2\n	allow =all\n", "user2", false},
	{"deny =all\n", "user", false},
	{"allow =user\n", "user", true},
	{"allow =all\n", "user", true},
	{"", "user", true},
}

func Test_isAccess(t *testing.T) {
	for _, tst := range tests {
		cfg := newFromKong(rawConfig{Allow: parseSetting(tst.conf, "allow"), Deny: parseSetting(tst.conf, "deny")}, "")
		t.Run(tst.conf, func(t *testing.T) {
			if res := cfg.IsAccess(tst.user); res != tst.want {
				t.Errorf("on user \"%s\" get %v want %v", tst.user, res, tst.want)
			}
		})
	}
}

func parseSetting(value, key string) string {
	for _, line := range strings.Split(value, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1])
		}
	}
	return ""
}

func TestLoadWithKongReadsJSONAndEnvironment(t *testing.T) {
	configFile := t.TempDir() + "/bot.cfg"
	if err := os.WriteFile(configFile, []byte(`{"company":"Configured Company","allow":"all"}`), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TELEGRAM_BOT_TOKEN", "test-token")

	values, err := loadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if values.Company != "Configured Company" || values.Token != "test-token" {
		t.Fatalf("values = %***REMOVED***v, want configured company and environment token", values)
	}
}

func TestLoadWithKongReadsDotEnvFile(t *testing.T) {
	dir := t.TempDir()
	configFile := dir + "/bot.cfg"
	envFile := dir + "/.env"
	if err := os.WriteFile(configFile, []byte(`{"company":"Configured Company"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(envFile, []byte("TELEGRAM_BOT_TOKEN=dotenv-token\n"), 0600); err != nil {
		t.Fatal(err)
	}

	values, err := loadWithKong(configFile, envFile)
	if err != nil {
		t.Fatal(err)
	}
	if values.Token != "dotenv-token" {
		t.Fatalf("token = %q, want %q", values.Token, "dotenv-token")
	}
}

func TestKongConfigLoaderAllowsComments(t *testing.T) {
	configFile := t.TempDir() + "/bot.cfg"
	content := "// comment\n{\n  // company name\n  \"company\": \"Test Company\"\n}"
	if err := os.WriteFile(configFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	values, err := loadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if values.Company != "Test Company" {
		t.Fatalf("company = %q, want %q", values.Company, "Test Company")
	}
}

func TestConfigUpdateKeepsLastGoodConfigOnInvalidFile(t *testing.T) {
	configFile := t.TempDir() + "/bot.cfg"
	if err := os.WriteFile(configFile, []byte(`{"company":"Working Company","allow":"all"}`), 0600); err != nil {
		t.Fatal(err)
	}

	conf, err := NewConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configFile, []byte(`{"company":`), 0600); err != nil {
		t.Fatal(err)
	}
	conf.Update()

	if conf.Company != "Working Company" {
		t.Fatalf("company after invalid update = %q, want %q", conf.Company, "Working Company")
	}
}

func TestLoadAppliesCLIOverrides(t *testing.T) {
	configFile := t.TempDir() + "/bot.cfg"
	if err := os.WriteFile(configFile, []byte(`{"company":"Config Company"}`), 0600); err != nil {
		t.Fatal(err)
	}

	conf, options, err := Load([]string{"--config", configFile, "--token", "cli-token", "--verbose", "--debug"})
	if err != nil {
		t.Fatal(err)
	}
	if conf.Token != "cli-token" {
		t.Fatalf("token = %q, want CLI token", conf.Token)
	}
	if !options.Verbose || !options.Debug {
		t.Fatalf("options = %***REMOVED***v, want verbose and debug enabled", options)
	}
}

func TestLoadRequiresToken(t *testing.T) {
	configFile := t.TempDir() + "/bot.cfg"
	if err := os.WriteFile(configFile, []byte(`{"company":"Config Company"}`), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TELEGRAM_BOT_TOKEN", "")

	if _, _, err := Load([]string{"--config", configFile}); err == nil {
		t.Fatal("expected missing token error")
	}
}

// func Test_st1(t *testing.T) {
// 	st()
// }
