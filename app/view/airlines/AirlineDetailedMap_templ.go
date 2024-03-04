// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.598
package airline

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/FACorreiaa/Aviation-tracker/app/models"

func detailedMapContainer(data models.Airline) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_detailedMapContainer_3a45`,
		Function: `function __templ_detailedMapContainer_3a45(data){function createFeatureFromAirline(airline) {

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
            geometry: new ol.geom.Point(ol.proj.fromLonLat([airline.longitude, airline.latitude])),
        });


        feature.setStyle(iconStyle);

        return feature;
    }


   const vectorSource = new ol.source.Vector({
    features: [createFeatureFromAirline(data)],
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
      controls: [],
      view: new ol.View({
         center: ol.proj.fromLonLat([data.longitude, data.latitude]),
         zoom: 10,
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

         //show toolip
      const contentDiv = document.createElement('div');
      contentDiv.innerHTML = ` + "`" + `
            <strong>Airport:</strong> ${feature.get('airport')}<br>
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

    map.getView().on('change:resolution', function () {
            if (map.getView().getZoom() < 4) {
                vectorLayer.setVisible(false);
            } else {
                vectorLayer.setVisible(true);
            }
        });
}`,
		Call:       templ.SafeScript(`__templ_detailedMapContainer_3a45`, data),
		CallInline: templ.SafeScriptInline(`__templ_detailedMapContainer_3a45`, data),
	}
}

func AirlineDetailedMap(data models.Airline) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html><head><style scoped>\n\t\t\t\t.map {\n\t\t\t\t\twidth: 100%;\n\t\t\t\t\theight: 500px;\n\t\t\t\t}\n\t\t\t\ta.skiplink {\n\t\t\t\t\tposition: absolute;\n\t\t\t\t\tclip: rect(1px, 1px, 1px, 1px);\n\t\t\t\t\tpadding: 0;\n\t\t\t\t\tborder: 0;\n\t\t\t\t\theight: 1px;\n\t\t\t\t\twidth: 1px;\n\t\t\t\t\toverflow: hidden;\n\t\t\t\t}\n\t\t\t\ta.skiplink:focus {\n\t\t\t\t\tclip: auto;\n\t\t\t\t\theight: auto;\n\t\t\t\t\twidth: auto;\n\t\t\t\t\tbackground-color: #fff;\n\t\t\t\t\tpadding: 0.3em;\n\t\t\t\t}\n\t\t\t\t#map:focus {\n\t\t\t\t\toutline: #4A74A8 solid 0.15em;\n\t\t\t\t}\n\t\t\t</style></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, detailedMapContainer(data))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body onload=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 templ.ComponentScript = detailedMapContainer(data)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2.Call)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"mb-5 text-left\"><button id=\"zoom-out\" class=\"btn btn-ghost mr-5\">Zoom out</button> <button id=\"zoom-in\" class=\"btn btn-ghost\">Zoom in</button></div><div id=\"map\" class=\"map\" tabindex=\"0\"><button aria-describedby=\"popup\" data-tippy-content=\"popup\" id=\"popup\"></button></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
