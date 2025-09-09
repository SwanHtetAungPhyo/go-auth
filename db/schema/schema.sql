
CREATE TABLE IF NOT EXISTS goauth_user (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                           email VARCHAR(255) UNIQUE NOT NULL,
                                           hash_password TEXT NOT NULL,
                                           name VARCHAR(255),
                                           image TEXT,
                                           role_name varchar(60) not null  default 'USER',
                                           email_verified BOOLEAN DEFAULT FALSE,
                                           two_factor_enabled BOOLEAN DEFAULT FALSE,
                                           two_factor_secret TEXT,
                                           metadata JSONB DEFAULT '{}',
                                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE  TABLE  IF NOT EXISTS  goauth_audit_log(
  id serial primary key ,
  event_type varchar not null ,
  log_entry jsonb not null default  '{}',
  occurred_at timestamp default  now()
);
-- Create social accounts table
CREATE TABLE IF NOT EXISTS goauth_account (
                                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                              user_id UUID NOT NULL REFERENCES goauth_user(id) ON DELETE CASCADE,
                                              provider VARCHAR(50) NOT NULL,
                                              provider_id VARCHAR(255) NOT NULL,
                                              created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                              UNIQUE(provider, provider_id)
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS goauth_session (
                                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                              user_id UUID NOT NULL REFERENCES goauth_user(id) ON DELETE CASCADE,
                                              token TEXT UNIQUE NOT NULL,
                                              expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                              user_agent TEXT,
                                              ip_address INET,
                                              created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create password reset tokens table
CREATE TABLE IF NOT EXISTS goauth_password_reset (
                                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                     user_id UUID NOT NULL REFERENCES goauth_user(id) ON DELETE CASCADE,
                                                     token TEXT UNIQUE NOT NULL,
                                                     expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create email verification tokens table
CREATE TABLE IF NOT EXISTS goauth_email_verification (
                                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                         user_id UUID NOT NULL REFERENCES goauth_user(id) ON DELETE CASCADE,
                                                         token TEXT UNIQUE NOT NULL,
                                                         expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_goauth_user_email ON goauth_user(email);
CREATE INDEX IF NOT EXISTS idx_goauth_session_user_id ON goauth_session(user_id);
CREATE INDEX IF NOT EXISTS idx_goauth_session_token ON goauth_session(token);
CREATE INDEX IF NOT EXISTS idx_goauth_session_expires_at ON goauth_session(expires_at);
CREATE INDEX IF NOT EXISTS idx_goauth_account_user_id ON goauth_account(user_id);
CREATE INDEX IF NOT EXISTS idx_goauth_account_provider ON goauth_account(provider, provider_id);
