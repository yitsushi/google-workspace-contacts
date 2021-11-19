# Google Workspace Contacts for (Neo)Mutt


## Install

```
go install github.com/yitsushi/google-workspace-contacts@v1.0.0
```

## Usage

Move (save) your `credentials.json` from Google Cloud Console under
`${HOME}/.config/google-workspace-contacts/credentials.json`

```
❯ google-workspace-contacts -h
Usage of google-workspace-contacts:
  -output-file string
        Output file, default to stdout (default "-")
  -v    Verbose output
```

You can generate your aliases list by redirecting the output (if you want to filter the output with other tools) or specify an output file.

```
❯ google-workspace-contacts -output-file ~/.mutt/weaveworks/ww_aliases
```

On first run it will ask you to open a URL and copy the authorization code.

```
Go to the following link in your browser then copy back the authorization code:
https://accounts.google.com/o/oauth2/auth?very-long-long-long-url
Token: token you get in your browser after granted permissions to the application
```

## Limitations

### Missing Names / Nicknames

It works, but Google API is a piece of ^&-%. It does not return names for a lot of users.

References:
* https://stackoverflow.com/questions/65920897/google-people-api-listdirectorypeople-is-not-returning-names-for-some-users-in-t
* https://stackoverflow.com/questions/65801879/google-people-api-listdirectorypeople-is-not-returning-names-for-any-users-in-th

As a workaround
* If Nickname is empty, the system generates Nickname from the email address
  (everything before `@`)
* If Name is empty, the system generates a Name from the email address by
  replacing all `.` with a space character and converting first letters to
  uppercase. (everything before `@`)
