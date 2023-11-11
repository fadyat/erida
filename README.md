### Erida: Simplifying Internal Cluster Communication

**Erida** is a straightforward SMTP relay server designed for sending internal cluster emails to an authenticated SMTP
server with seamless Slack integration.

### Email Address Flexibility

Erida supports a variety of email addresses, including common ones like `avfadeev@gmail.com` and those associated with
messaging services, particularly Slack.

For Slack integration, you can use addresses following this syntax:

- User-specific Slack address: `personal.<username>@slack`
- Channel-specific Slack address: `channel.<channelname>@slack`

It's important to note that both the `username` and `channelname` are case-insensitive, and the bot must have the
necessary
permissions to access the specified Slack channels.

### Configuration Made Easy

Configuring Erida is a breeze. All you need to do is set the following environment variables:

- `SMTP_HOST`: SMTP server host
- `SMTP_PORT`: SMTP server port
- `SMTP_USER`: SMTP server username
- `SMTP_PASS`: SMTP server password
- `SLACK_TOKEN`: Slack bot token
- `SMTP_TLS` (Optional, default: true): Enable or disable Start TLS usage

For additional variables, refer to the [configuration file](internal/config.go).

### Getting Started

If you're new to configuring the bot, check out the step-by-step guide
at [Slack Quickstart](https://api.slack.com/start/quickstart).

Ensure that the bot has the necessary permissions, specifically `chat:write`.

### Example Usage

Let's walk through an example. Assume that an external SMTP server is configured to send emails to **Erida** with the
following addresses: `personal.fadyat@slack`, `channel.general@slack`, and `avfadeev@gmail.com`.

The message will be seamlessly delivered to the Slack channel `#general` and the Slack user `@fadyat`, as well as to the
email address `avfadeev@gmail.com`.

**Erida** simplifies internal communication, bridging the gap between email and Slack effortlessly.

