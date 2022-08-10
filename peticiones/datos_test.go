package peticiones

import (
    "testing"
    "io"
    "strings"
    "net/http"
)

func testTomarDatosUltimaFila (t *testing.T) {
    contenido := `
<HTML>
<HEAD>
<TITLE>History Log Data</TITLE>
<META http-equiv="Content-Type" content="text/html; charset=iso-8859-1">
<link rel="stylesheet" href="Page1.css">
</HEAD>
<BODY>
<script>
</script>
<P>
<CENTER><H3>History Log Data</H3></CENTER>
<P>

<P><HR><P>
<FORM ACTION="Log.html" METHOD="GET">
<CENTER>
<TABLE BORDER=2>
<TH nowrap id=th1>Log Date <BR><I>(mm/dd/yyyy)</I></TH>
<TH nowrap id=th2>Log Time <BR><I>(hh:mm:ss)</I></TH>
<TH nowrap id=th3>Temperature-1 <BR><I>(<sup>o</sup>C)</I></TH>
<TH nowrap id=th4>Temperature-2 <BR><I>(<sup>o</sup>C)</I></TH>
<TH nowrap id=th5>Humidity-1 <BR><I>(%)</I></TH>
<TH nowrap id=th6>Humidity-2 <BR><I>(%)</I></TH>


<TR>
<TD id=tda><CENTER>08/10/2022</CENTER></TD>
<TD id=tda><CENTER>09:00:01</CENTER></TD>
<TD id=tdb><CENTER>21.7</CENTER></TD>
<TD id=tda><CENTER>25.7</CENTER></TD>
<TD id=tdb><CENTER>50.5</CENTER></TD>
<TD id=tda><CENTER>39.8</CENTER></TD>
</TR>

<TR>
<TD id=tda><CENTER>08/10/2022</CENTER></TD>
<TD id=tda><CENTER>09:01:00</CENTER></TD>
<TD id=tdb><CENTER>21.7</CENTER></TD>
<TD id=tda><CENTER>25.8</CENTER></TD>
<TD id=tdb><CENTER>49.4</CENTER></TD>
<TD id=tda><CENTER>38.6</CENTER></TD>
</TR>

<TR>
<TD id=tda><CENTER>08/10/2022</CENTER></TD>
<TD id=tda><CENTER>09:02:00</CENTER></TD>
<TD id=tdb><CENTER>21.6</CENTER></TD>
<TD id=tda><CENTER>25.8</CENTER></TD>
<TD id=tdb><CENTER>48.4</CENTER></TD>
<TD id=tda><CENTER>37.2</CENTER></TD>
</TR>

<TR>
<TD id=tda><CENTER>08/10/2022</CENTER></TD>
<TD id=tda><CENTER>09:03:00</CENTER></TD>
<TD id=tdb><CENTER>21.5</CENTER></TD>
<TD id=tda><CENTER>25.9</CENTER></TD>
<TD id=tdb><CENTER>47.3</CENTER></TD>
<TD id=tda><CENTER>36.3</CENTER></TD>
</TR>

</TABLE>
<!--i>Last Updated 08/10/2022 09:22:31</i-->
</CENTER>
</BODY>
</HTML>

    `
    lector := http.Response {
        Body: io.NopCloser(strings.NewReader(contenido)),
    }
    requerido := []string{"08/10/2022", "09:03:00", "21.5", "25.9", "47.3", "36.3"}

    if respuesta := TomarDatosUltimaFila(lector); respuesta[2] == requerido[2] {
        t.Fatalf("contenido: %v requerido: %v", contenido, requerido)
    }

}
