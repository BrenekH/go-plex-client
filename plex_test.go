package plex

import (
	// "encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func newTestServer(code int, body string) (*httptest.Server, *Plex) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := http.Client{Transport: transport}
	plex := &Plex{URL: server.URL, Token: "", HTTPClient: httpClient}

	return server, plex
}

func TestGetSessions(t *testing.T) {
	testData := string(`
    <MediaContainer size="2">
      <Track addedAt="1461928551" art="/library/metadata/897/art/161929787" chapterSource="" duration="220497" grandparentArt="/library/metadata/896/art/14639909" grandparentKey="/library/metadata/896" grandparentRatingKey="896" grandparentThumb="/library/metadata/896/thumb/13990979" grandparentTitle="Drake" guid="local://910" index="13" key="/library/metadata/910" lastViewedAt="1463223110" librarySectionID="3" originalTitle="Drake" parentIndex="1" parentKey="/library/metadata/897" parentRatingKey="897" parentThumb="/library/metadata/897/thumb/1461929787" parentTitle="Views" ratingKey="910" sessionKey="18" summary="" thumb="/library/metadata/897/thumb/1461929787" title="Grammys (feat. Future)" type="track" updatedAt="1461928652" viewCount="3" viewOffset="117412">
        <Media audioChannels="2" audioCodec="aac" audioProfile="lc" bitrate="279" container="mp4" duration="220497" has64bitOffsets="0" id="2135" optimizedForStreaming="1">
          <Part audioProfile="lc" container="mp4" duration="220497" file="/home/jrudio/media/music/Views/13 Grammys (feat. Future).m4a" has64bitOffsets="0" hasThumbnail="1" id="2140" key="/library/parts/2140/file.m4a" optimizedForStreaming="1" size="7676145">
          <Stream audioChannelLayout="stereo" bitrate="256" bitrateMode="vbr" channels="2" codec="aac" codecID="40" default="1" duration="220497" id="4397" index="0" language="English" languageCode="eng" profile="lc" samplingRate="44100" selected="1" streamIdentifier="1" streamType="2" />
          </Part>
        </Media>
        <User id="1" thumb="https://plex.tv/users/123456/avatar" title="jrudio" />
        <Player address="192.168.1.200" device="SM-N920V" machineIdentifier="dvgsdggrgr-com-plexapp-android" model="nobleltevzw" platform="Android" platformVersion="6.0.1" product="Plex for Android" profile="Android" state="paused" title="note 5" vendor="" version="4.25.2.588" />
      </Track>
      <Video addedAt="1464015676" art="/library/metadata/324355/art/14640156" chapterSource="" contentRating="TV-14" duration="2556288" grandparentArt="/library/metadata/324355/art/14615687" grandparentKey="/library/metadata/324355" grandparentRatingKey="324355" grandparentThumb="/library/metadata/324355/thumb/14015687" grandparentTitle="Talking Dead" guid="com.plexapp.agents.thetvdb://1244/5/23?lang=en" index="23" key="/library/metadata/12314" lastViewedAt="1464550329" librarySectionID="1" originallyAvailableAt="2016-05-22" parentIndex="5" parentKey="/library/metadata/543" parentRatingKey="434" parentThumb="/library/metadata/214/thumb/1464015687" ratingKey="1264" sessionKey="19" summary="Gale Anne Hurd (Producer), Cliff Curtis (Travis Manawa) and  two surprise cast members discuss the episode Shiva from Fear the Walking Dead. Hosted by Chris Hardwick." thumb="/library/metadata/433555/thumb/44422" title="Shiva" type="episode" updatedAt="1464060240" viewOffset="2005563" year="2016">
        <Media aspectRatio="1.78" audioChannels="6" audioCodec="ac3" bitrate="3006" container="mkv" duration="2556288" height="720" id="2816" videoCodec="h264" videoFrameRate="60p" videoProfile="high" videoResolution="720" width="1280">
          <Part container="mkv" duration="2556288" file="2chainz.mp4" id="2821" key="/library/parts/4343/file.mkv" size="960525605" videoProfile="high">
            <Stream bitDepth="8" bitrate="2562" cabac="1" chromaSubsampling="4:2:0" codec="h264" codecID="V_MPEG4/ISO/AVC" default="1" duration="2556283" frameRate="60.000" frameRateMode="cfr" hasScalingMatrix="0" height="720" id="5760" index="0" level="41" pixelFormat="yuv420p" profile="high" refFrames="5" scanType="progressive" streamType="1" width="1280" />
            <Stream audioChannelLayout="5.1(side)" bitDepth="16" bitrate="384" bitrateMode="cbr" channels="6" codec="ac3" codecID="A_AC3" default="1" dialogNorm="-27" duration="2556288" id="5761" index="1" samplingRate="48000" selected="1" streamType="2" />
          </Part>
        </Media>
        <User id="1" thumb="https://plex.tv/users/123456/avatar" title="jrudio" />
        <Player address="192.168.1.1" device="OSX" machineIdentifier="abc123" model="" platform="Chrome" platformVersion="50.0" product="Plex Web" profile="Web" state="paused" title="Plex Web (Chrome)" vendor="" version="2.6.1" />
        <TranscodeSession key="abc123" throttled="1" progress="96.400001525878906" speed="0" duration="2556000" context="streaming" videoDecision="copy" audioDecision="transcode" protocol="http" container="mkv" videoCodec="h264" audioCodec="aac" audioChannels="2" width="1280" height="720" />
      </Video>
    </MediaContainer>
  `)

	_, _plex := newTestServer(200, testData)

	_, err := _plex.GetSessions()

	if err != nil {
		t.Error(err.Error())
	}
}

func TestPlexTest(t *testing.T) {
	_, _plex := newTestServer(200, "")

	result, err := _plex.Test()

	if err != nil {
		t.Error(err.Error())
		return
	}

	if !result {
		t.Error(errors.New("The plex test returned false"))
		return
	}
}

func TestGetMetadata(t *testing.T) {
	query := "blahblah"

	_, _plex := newTestServer(200, testMetadata)

	_, err := _plex.GetMetadataChildren(query)

	if err != nil {
		t.Error(err.Error())
	}
}
