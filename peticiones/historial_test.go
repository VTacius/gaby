package peticiones

import (
    "testing"
    "io"
    "strings"
    "net/http"
)

func TestEncontrarEnlaceListaHistorial(t *testing.T){
    contenido := `
<HTML>
<HEAD>
<TITLE>History Log Index Page </TITLE>
<link rel="stylesheet" href="Page1.css">
<META http-equiv="Content-Type" content="text/html; charset=iso-8859-1">
</HEAD>
<BODY id=bd1>
<script>
</script>
<p>
<CENTER><H3>History Log Index Page</H3></CENTER>
<P>

<HR><P><UL>
<CENTER>

<NOBR><LI><A HREF="XAAAAAAKOCBBALog11.html">From 08/09/2022 01:30:00</A></NOBR>

<NOBR><LI><A HREF="XAAAAAAJACBBALog11.html">From 08/09/2022 01:00:00 To 08/09/2022 01:29:00</A></NOBR>

<NOBR><LI><A HREF="XAAAAAAHCCBBALog11.html">From 08/09/2022 00:30:00 To 08/09/2022 00:59:00</A></NOBR>

<NOBR><LI><A HREF="XAAAAAAFECBBALog11.html">From 08/09/2022 00:00:00 To 08/09/2022 00:29:00</A></NOBR>

<NOBR><LI><A HREF="XAAAAAADGCBBALog11.html">From 08/08/2022 23:30:00 To 08/08/2022 23:59:00</A></NOBR>

<NOBR><LI><A HREF="XAAAAAABICBBALog11.html">From 08/08/2022 23:00:00 To 08/08/2022 23:29:00</A></NOBR>

<NOBR><LI><A HREF="XAAAAAAABCBBALog11.html">From 08/08/2022 22:37:01 To 08/08/2022 22:59:00</A></NOBR>

</CENTER>
<!--i>Last Updated 08/10/2022 08:15:47</i-->
</UL>
</BODY>
</HTML>
`
    lector := http.Response {
        Body: io.NopCloser(strings.NewReader(contenido)),
    }
    requerido := "XAAAAAAKOCBBALog11.html"
    if respuesta := EncontrarEnlaceListaHistorial(lector); respuesta != requerido {
        t.Fatalf("contenido: %s requerido: %s", contenido, requerido)
    }
}
