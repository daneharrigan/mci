package notifier

const Template = `
	<h1 style="font: bold 24px Helvetica; color: #111;">MCI: New Release This Week</h1>
	<table style="padding: 0; margin: 0; border: none;">
	{{range .}}
	  <tr>
		<td style="margin: 0; padding: 20px 0; border: none;"><img src="{{.Thumbnail}}" style="width: auto; height: 200px"></td>
		<td style="margin: 0; padding: 20px 0; border: none;">
		  <p style="padding: 0; margin: 0; border: none; font: normal 14px Helvetica; color: #666;"><strong>{{.Name}}</strong></p>
		  <p style="padding: 0; margin: 0; border: none; font: normal 14px Helvetica; color: #666;">{{series .}}</p>
	  </tr>
	{{end}}
	</table>
	`
