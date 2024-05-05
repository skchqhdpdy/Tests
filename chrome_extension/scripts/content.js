setTimeout(() => {
    const ui_vertical_buttons = document.getElementsByClassName("ui vertical buttons")[0]
    const beatmapID = "..."
    const setData = {"SetID": "..."}

    //audio Tag
    const audio_tag = document.createElement("audio")
    audio_tag.id = "BeatmapAudio"; audio_tag.innerText = "Your browser does not support the audio element."
    const audio_source_tag = document.createElement("source")
    audio_source_tag.id = "song_source"; audio_source_tag.src = `https://b.redstar.moe/preview/${setData['SetID']}`; audio_source_tag.type = "audio/mp3";
    audio_tag.appendChild(audio_source_tag)
    document.body.appendChild(audio_tag)

    //preview1
    const preview1 = document.createElement("a")
    preview1.onclick = togglePlay(); preview1.class = "ui blue labeled icon button";
    var i = document.createElement("i")
    i.id = "imageplay"; i.class = "play icon"; preview1.appendChild(i)
    var span = document.createElement("span")
    span.id = "PlayState"; span.innerText = "Play Song"; preview1.appendChild(span)
    ui_vertical_buttons.appendChild(preview1)

    //preview2
    const preview2 = document.createElement("a")
    preview2.class = "ui blue labeled icon button";
    var i = document.createElement("i")
    i.id = "imageplay2"; i.class = "play icon"; preview2.appendChild(i)
    var span = document.createElement("span")
    span.id = "current_time"; span.style.color = "red"; span.innerText = "&#11014 wait a sec..."; preview2.appendChild(span)
    ui_vertical_buttons.appendChild(preview2)

    //osu!direct
    const direct = document.createElement("a")
    direct.href = `osu://b/${beatmapID}`; direct.class = "ui pink labeled icon button"; direct.innerText = "osu!direct"
    var i = document.createElement("i")
    i.class = "download icon"; direct.appendChild(i)
    ui_vertical_buttons.appendChild(direct)

    //View On Bancho
    const VOB = document.createElement("a")
    VOB.target = "_blank"; VOB.href = `https://osu.ppy.sh/b/${beatmapID}`; VOB.class = "ui pink labeled icon button"; VOB.innerText = "View On Bancho"
    var i = document.createElement("i")
    i.class = "play icon"; VOB.appendChild(i)
    ui_vertical_buttons.appendChild(VOB)

    console.log(window.beatmapID)
}, 500);

setTimeout(() => {
    //script
    //for playing and stopping bmap music
    var audio = document.getElementById("BeatmapAudio");
    console.log(audio)

    audio.addEventListener('loadedmetadata', () => {
        document.getElementById("current_time").innerText = `0.0 / ${Math.round(audio.duration * 10) / 10}`
        document.getElementById("current_time").style.color = ""
    });
    // 현재 오디오 플레이 정보
    audio.addEventListener("timeupdate", function(){
        document.getElementById("current_time").innerText = `${Math.round(audio.currentTime * 10) / 10} / ${Math.round(audio.duration * 10) / 10}`
    });
    audio.addEventListener("ended", function(){
        document.getElementById("current_time").innerText = `0 / ${Math.round(audio.duration * 10) / 10}`
    });
    // Realistik Server Code
    var PlayText = 'Play Song'
    var PauseText = 'Pause Song'
    var isPlaying = false;

    audio.volume = 0.2;

    function togglePlay() {
        if (isPlaying) {
        audio.pause()
        document.getElementById("PlayState").innerHTML = PlayText;
        document.getElementById("imageplay").classList.remove('pause');
        document.getElementById("imageplay").classList.add('play');

        document.getElementById("imageplay2").classList.remove('pause');
        document.getElementById("imageplay2").classList.add('play');
        } else {
        console.log(`Audio Volume = ${audio.volume}`)

        audio.play();
        document.getElementById("PlayState").innerHTML = PauseText;
        document.getElementById("imageplay").classList.remove('play');
        document.getElementById("imageplay").classList.add('pause');

        document.getElementById("imageplay2").classList.remove('play');
        document.getElementById("imageplay2").classList.add('pause');
        }
    };
    audio.onplaying = function() {
        isPlaying = true;
    };
    audio.onpause = function() {
        isPlaying = false;
    };

    audio.addEventListener("ended", function(){
        audio.currentTime = 0;
        document.getElementById("PlayState").innerHTML = PlayText;
        document.getElementById("imageplay").classList.remove('pause');
        document.getElementById("imageplay").classList.add('play');

        document.getElementById("imageplay2").classList.remove('pause');
        document.getElementById("imageplay2").classList.add('play');
    });
}, 1000);

/*
<!-- 원래 hanayo 소스 -->
                  <a onclick="togglePlay()" class="ui blue labeled icon button">
                    <i id="imageplay" class="play icon"></i>
                    <span id="PlayState">
                      Play Song
                    </span>
                  </a>
                  <a class="ui blue labeled icon button">
                    <i id="imageplay2" class="play icon"></i>
                    <span id="current_time" style="color: red;">								
                      &#11014 wait a sec...
                    </span>
                  </a>
                  
                  <a href="osu://b/{{ .Beatmap.ID }}" class="ui pink labeled icon button"><i class="download icon"></i>{{ $.T "osu!direct" }}</a>
                  <a href="https://redstar.moe/d/{{ .Beatmapset.ID }}" class="ui green labeled icon button"><i class="download icon"></i>{{ $.T "download" }}</a>
                  {{ if .Beatmapset.HasVideo }}
                    <a href="https://redstar.moe/d/{{ .Beatmapset.ID }}?novideo" class="ui gray labeled icon button"><i class="download icon"></i>{{ $.T "download (no video)" }}</a>
                    {{ end }}
                  {{ if has $.Context.User.Privileges 256 }}
                    <a href="https://old.redstar.moe/index.php?p=124&bsid={{ .Beatmapset.ID }}" class="ui violet labeled icon button"><i class="thumbs up icon"></i>{{ $.T "Rank in RAP"}}</a>
                    {{ end }}
  
  
                    {{ if has $.Context.User.Privileges 267 }}
                    <a target="_blank" href="https://admin.redstar.moe/rank/{{ .Beatmap.ID }}" class="ui yellow labeled icon button"><i class="folder open icon"></i>{{ $.T "Manage Beatmap"}}</a>
                    {{ end }}
                    
                    <a target="_blank" href="https://osu.ppy.sh/b/{{ .Beatmap.ID }}" class="ui pink labeled icon button"><i class="play icon"></i>{{ $.T "View On Bancho" }}</a>
  
                    <a target="_blank" href="https://preview.nerinyan.moe/#{{ .Beatmap.ID }}" class="ui violet labeled icon button"><i class="play icon"></i>{{ $.T "Preview Beatmap" }}</a>
                    <a target="_blank" href="https://osu.pages.dev/preview#{{ .Beatmap.ID }}" class="ui violet labeled icon button"><i class="play icon"></i>{{ $.T "Preview Beatmap2" }}</a>
                    <!-- 원래 hanayo 소스 끝 -->

<audio id="BeatmapAudio">
    <source id="song_source" src="https://b.redstar.moe/preview/{{ .Beatmapset.ID }}.mp3" type="audio/mp3">
    Your browser does not support the audio element.
</audio>

*/