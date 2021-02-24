# certtools

`certtools` is created to manage certificates store in Windows Nano Containers. Where `certutil.exe` and `powershell Import-PfxCertificate` are not available.

# Commands

## Import

```
certtools.exe import [command options] <path/to/pfx>
   --file value, -f value      path to pfx
   --password value, -p value  password to pfx
```

the first arg and -f are identical

Example:

```
C:\>certtools.exe import test.pfx
```

## List

Example:

```
C:\>certtools.exe ls
a17aa4e9afc16f8cc15864703aa8186e58daddbe test 
```

## Remove

```
certtools.exe rm [command options] <thumbprint>
   --thumbprint value, -t value  thumbprint of the certificate to be deleted
```

the first arg and -t are identical

Example:

```
C:\>certtools.exe rm a17aa4e9afc16f8cc15864703aa8186e58daddbe
```

