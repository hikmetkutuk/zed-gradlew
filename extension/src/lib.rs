use zed_extension_api::{self as zed, LanguageServerId, Result};

struct GradlewExtension;

impl zed::Extension for GradlewExtension {
    fn new() -> Self {
        GradlewExtension
    }

    fn language_server_command(
        &mut self,
        _language_server_id: &LanguageServerId,
        worktree: &zed::Worktree,
    ) -> Result<zed::Command> {
        let binary = worktree
            .which("zed-gradlew")
            .unwrap_or_else(|| "/Users/hikmetkutuk/.local/bin/zed-gradlew".to_string());

        Ok(zed::Command {
            command: binary,
            args: vec![],
            env: vec![],
        })
    }
}

zed::register_extension!(GradlewExtension);
