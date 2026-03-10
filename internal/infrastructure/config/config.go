package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type DatabaseConfig struct {
	DSN string
}

func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		DSN: os.Getenv("DB_DSN"),
	}
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

func LoadRedisConfig() RedisConfig {
	db := 0
	if val := os.Getenv("REDIS_DB"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			db = i
		}
	}

	appName := envStringOrDefault("APP_NAME", "laravel")
	defaultRedisPrefix := slugify(appName) + "_database_"
	redisPrefix := envStringOrDefault("REDIS_PREFIX", defaultRedisPrefix)

	return RedisConfig{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
		Prefix:   redisPrefix,
	}
}

type AsynqRedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func LoadAsynqRedisConfig() AsynqRedisConfig {
	db := 1
	if val := os.Getenv("ASYNQ_REDIS_DB"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			db = i
		}
	}
	return AsynqRedisConfig{
		Addr:     os.Getenv("ASYNQ_REDIS_HOST") + ":" + os.Getenv("ASYNQ_REDIS_PORT"),
		Password: os.Getenv("ASYNQ_REDIS_PASSWORD"),
		DB:       db,
	}
}

func envStringOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func envIntOrDefault(key, defaultValue string) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return 1
}

// slugify simulates Laravel's Str::slug with an underscore separator
func slugify(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "_")
	return strings.Trim(s, "_")
}

type CacheConfig struct {
    Prefix   string
    Prefixes CachePrefixes
    Tags     CacheTags
    Expiry   CacheExpiry
}

type CachePrefixes struct {
    Till                    string
    PosAuthUser             string
    PosAuthToken            string
    Store                   string
    User                    string
    DeviceRemember          string
    Device                  string
    Product                 string
    DuplicateCheck          string
    CustomerOTP             string
    SPOSAuth                string
    DraftBill               string
    DraftBillTags           string
    ProductMolecule         string
    OngcPrescription        string
    DraftPaymentDelete      string
    ScreenLink              string
    ProductImage            string
    FirebaseJWTToken        string
    OngcLoginAccessToken    string
    OngcLoginRefreshToken   string
}

type CacheTags struct {
    Till            string
    Product         string
    ProductMolecule string
    OngcPrescription string
    DraftBill       string
    ScreenLink      string
    ProductImage    string
}

type CacheExpiry struct {
    TillHours                    int
    PosAuthHours                 int
    DuplicateCheckSeconds        int
    OngcPrescriptionHours        int
    DraftPaymentDeleteSeconds    int
    ScreenLinkMinutes            int
    ProductImageHours            int
    OngcLoginRefreshTokenHours   int
}

func LoadCacheConfig() CacheConfig {
    appName := envStringOrDefault("APP_NAME", "laravel")
    defaultPrefix := slugify(appName) + "_cache_:"
    prefix := envStringOrDefault("CACHE_PREFIX", defaultPrefix)
    
    return CacheConfig{
        Prefix:   prefix,
        Prefixes: loadCachePrefixes(),
        Tags:     loadCacheTags(),
        Expiry:   loadCacheExpiry(),
    }
}

func loadCachePrefixes() CachePrefixes {
    return CachePrefixes{
        Till:                   envStringOrDefault("TILL_CACHE_PREFIX", "till_cache_data_"),
        PosAuthUser:            envStringOrDefault("POS_AUTH_USER_CACHE_KEY_PREFIX", "pos_auth_user_cache_"),
        PosAuthToken:           envStringOrDefault("POS_AUTH_CACHE_KEY_PREFIX", "pos_auth_token_"),
        Store:                  envStringOrDefault("STORE_CACHE_KEY", "STORE_CACHE_KEY_"),
        User:                   envStringOrDefault("USER_CACHE_KEY", "USER_CACHE_KEY_"),
        DeviceRemember:         envStringOrDefault("DEVICE_REMEMBER_CACHE_KEY", "device_remember"),
        Device:                 envStringOrDefault("DEVICE_CACHE_KEY", "device:"),
        Product:                envStringOrDefault("PRODUCT_CACHE_PREFIX", "product:"),
        DuplicateCheck:         envStringOrDefault("DUPLICATE_CHECK_CACHE_KEY", "duplicate_check_cache_key_"),
        CustomerOTP:            envStringOrDefault("CUSTOMER_OTP_CACHE_KEY", "CUSTOMER_OTP_CACHE_KEY_"),
        SPOSAuth:               envStringOrDefault("SPOS_AUTH_CACHE_KEY", "spos_auth_token_"),
        DraftBill:              envStringOrDefault("PREFIX_DRAFT_BILL_CACHE", "draft_invoices_"),
        DraftBillTags:          envStringOrDefault("PREFIX_DRAFT_BILL_CACHE_TAGS", "draft_invoices_tags_"),
        ProductMolecule:        envStringOrDefault("PRODUCT_MOLECULE_CACHE_PREFIX", "product_molecule_cache_prefix_"),
        OngcPrescription:       envStringOrDefault("ONGC_PRESCRIPTION_CACHE_PREFIX", "ongc_prescription_unique_id_"),
        DraftPaymentDelete:     envStringOrDefault("DRAFT_PAYMENT_DELETE_CACHE_KEY", "draft_payment_delete_"),
        ScreenLink:             envStringOrDefault("SCREEN_LINK_CACHE_", "SCREEN_LINK_CACHE_KEY_"),
        ProductImage:           envStringOrDefault("PRODUCT_IMAGE_CACHE_PREFIX", "product_image_"),
        FirebaseJWTToken:       envStringOrDefault("FIREBASE_JWT_TOKEN_CACHE_KEY", "firebase_jwt_token"),
        OngcLoginAccessToken:   envStringOrDefault("ONGC_LOGIN_ACCESS_TOKEN_CACHE_PREFIX", "ongc_login_access_token"),
        OngcLoginRefreshToken:  envStringOrDefault("ONGC_LOGIN_REFRESH_TOKEN_CACHE_PREFIX", "ongc_login_refresh_token"),
    }
}

func loadCacheTags() CacheTags {
    return CacheTags{
        Till:             envStringOrDefault("TILL_CACHE_TAGS", "till_cache_tags"),
        Product:          envStringOrDefault("PRODUCT_CACHE_TAG", "product_cache"),
        ProductMolecule:  envStringOrDefault("PRODUCT_MOLECULE_CACHE_TAG", "product_molecule_cache_tag"),
        OngcPrescription: envStringOrDefault("ONGC_PRESCRIPTION_CACHE_TAG", "ongc_prescription"),
        DraftBill:        envStringOrDefault("PREFIX_DRAFT_BILL_CACHE_TAGS", "draft_invoices_tags_"),
        ScreenLink:       envStringOrDefault("SCREEN_LINK_CACHE_TAG", "SCREEN_LINK_CACHE"),
        ProductImage:     envStringOrDefault("PRODUCT_IMAGE_CACHE_TAG", "product_image_cache"),
    }
}

func loadCacheExpiry() CacheExpiry {
    return CacheExpiry{
        TillHours:                  envIntOrDefault("TILL_CACHE_EXPIRY", "12"),
        PosAuthHours:               envIntOrDefault("POS_AUTH_CACHE_DURATION_HOURS", "12"),
        DuplicateCheckSeconds:      envIntOrDefault("DUPLICATE_REQUEST_CHECK_EXPIRY_IN_SECONDS", "10"),
        OngcPrescriptionHours:      envIntOrDefault("ONGC_PRESCRIPTION_CACHE_EXPIRY_HOURS", "24"),
        DraftPaymentDeleteSeconds:  envIntOrDefault("DRAFT_PAYMENT_DELETE_CACHE_EXPIRY_SECONDS", "5"),
        ScreenLinkMinutes:          envIntOrDefault("SCREEN_LINK_CACHE_EXPIRY_MINUTES", "10"),
        ProductImageHours:          envIntOrDefault("PRODUCT_IMAGE_CACHE_EXPIRY_HOURS", "18"),
        OngcLoginRefreshTokenHours: envIntOrDefault("ONGC_LOGIN_REFRESH_TOKEN_CACHE_EXPIRY_HOURS", "720"),
    }
}