// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.598
package locations

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"
	"github.com/FACorreiaa/Aviation-tracker/app/models"
)

func countryMapContainer(data []models.Country) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_countryMapContainer_2eee`,
		Function: `function __templ_countryMapContainer_2eee(data){//control selector
        const rangeInput = document.querySelector('.range');
        const updatePointsOnMap = () => {
            const selectedValue = parseInt(rangeInput.value, 10);
            // Logic to update the number of points on the map based on the selected value
            const filteredData = data.slice(0, selectedValue);
            // Clear existing features
            vectorSource.clear();
            // Add new features based on the filtered data
            vectorSource.addFeatures(filteredData.map(city => createFeatureFromCity(city)));
        };

    // Add event listener for input change
    rangeInput.addEventListener('input', updatePointsOnMap);

  function createFeatureFromCountry(country) {

        const iconStyle = new ol.style.Style({
            image: new ol.style.Icon({
                anchor: [1, 46],
                anchorXUnits: 'fraction',
                anchorYUnits: 'pixels',
                src: '../../static/icons/marker.png',
                scale: 0.5,
            }),
        });

        const feature = new ol.Feature({
            geometry: new ol.geom.Point(ol.proj.fromLonLat([country.longitude, country.latitude])),
            country: country.country_name,
            capital: country.capital,
            continent: country.continent,
            currency: country.currency_name,
            population: country.population,
            style: iconStyle,
        });

        feature.setStyle(iconStyle);

        return feature;
    }


    const vectorSource = new ol.source.Vector({
        features: data.map(country => createFeatureFromCountry(country)),
    });

   const vectorLayer = new ol.layer.Vector({
      source: vectorSource,
   });

  const tileLayer = new ol.layer.Tile({
            source: new ol.source.OSM(),
         })

   const map = new ol.Map({
      layers: [tileLayer, vectorLayer],
      target: document.getElementById('map'),
      view: new ol.View({
         center: [0, 0],
         zoom: 1,
      }),
   });

   const element = document.getElementById('popup');
   const popup = new ol.Overlay({
      element: element,
      positioning: 'bottom-center',
      stopEvent: false,
   });
   map.addOverlay(popup);

   let popover;

   function disposePopover() {
      if (popover) {
         popover.dispose();
         popover = undefined;
      }
   }

   const tippyButton = document.getElementById('popup');
    tippy(tippyButton, {
      content: document.createElement('div'),
      interactive: true,
      trigger: 'click',
      placement: 'top',
      animation: 'scale'  ,
      theme: 'translucent'
    });

   map.on('click', function (evt) {
      const feature = map.forEachFeatureAtPixel(evt.pixel, function (feature) {
         return feature;
      });
      disposePopover();
      if (!feature) {
         return;
      }
      popup.setPosition(evt.coordinate);

      const contentDiv = document.createElement('div');
      contentDiv.innerHTML = ` + "`" + `
            <strong>Country:</strong> ${feature.get('country')}<br>
            <strong>Capital:</strong> ${feature.get('capital')}<br>
            <strong>Continent:</strong> ${feature.get('continent')}<br>
            <strong>Currency:</strong> ${feature.get('currency')}<br>
            <strong>Population:</strong> ${feature.get('population')}<br>

    ` + "`" + `;

      tippyButton._tippy.setContent(contentDiv);
      tippyButton._tippy.show();
   });



   map.on('pointermove', function (e) {
      const pixel = map.getEventPixel(e.originalEvent);
      const hit = map.hasFeatureAtPixel(pixel);
      map.getTarget().style.cursor = hit ? "pointer" : "";
   });


   map.on('movestart', disposePopover);


   document.getElementById('zoom-out').onclick = function () {
      const view = map.getView();
      const zoom = view.getZoom();
      view.setZoom(zoom - 1);
   };

   document.getElementById('zoom-in').onclick = function () {
      const view = map.getView();
      const zoom = view.getZoom();
      view.setZoom(zoom + 1);
   };

   map.on('dblclick', event => {
       // get the feature you clicked
       const feature = map.forEachFeatureAtPixel(event.pixel, (feature) => {
        return feature
       })
       if(feature instanceof ol.Feature){
         // Fit the feature geometry or extent based on the given map
         map.getView().fit(feature.getGeometry())
         // map.getView().fit(feature.getGeometry().getExtent())
       }
   })
}`,
		Call:       templ.SafeScript(`__templ_countryMapContainer_2eee`, data),
		CallInline: templ.SafeScriptInline(`__templ_countryMapContainer_2eee`, data),
	}
}

func CountryMap(data []models.Country) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html><head><style scoped>\n\t\t\t\t.map {\n\t\t\t\t\twidth: 100%;\n\t\t\t\t\theight: 700px;\t\t\t\t}\n\t\t\t\t#map:focus {\n\t\t\t\t\toutline: #4A74A8 solid 0.15em;\n\t\t\t\t}\n\t\t\t</style></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, countryMapContainer(data))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body onload=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 templ.ComponentScript = countryMapContainer(data)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div id=\"map\" class=\"map\" tabindex=\"0\"><button aria-describedby=\"popup\" data-tippy-content=\"popup\" id=\"popup\"></button></div><div class=\"mt-2 text-center\"><button id=\"zoom-out\" class=\"btn btn-secondary\">Zoom out</button> <button id=\"zoom-in\" class=\"btn btn-secondary\">Zoom in</button></div><div class=\"w-full form-control\"><div class=\"label\"><span class=\"label-text font-semi-bold text-xs badge-xs badge pb-0\">Display number of markers</span></div><div class=\"flex items-center mt-2\"><span class=\"text-xs\">0</span> <input type=\"range\" min=\"0\" max=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", len(data))))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("%d", len(data))))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"range range-xs p-10 mr-2 ml-2\" name=\"rangeValue\"> <span class=\"text-xs\" id=\"rangeValue\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", len(data)))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/view/locations/CountryMap.templ`, Line: 192, Col: 99}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
