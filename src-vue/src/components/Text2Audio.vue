<template>
  <div class="t2a-container">
    <div class="grid base-settings">
      <h3>平台选择</h3>
      <!-- <span class="label">智能引擎选择：</span> -->
      <el-radio-group v-model="input.api">
        <el-radio-button
          v-for="(item, index) in apiList"
          :key="index"
          :label="item.value"
        >{{item.label}}</el-radio-button>
      </el-radio-group>

      <h3>基础设置</h3>
      <div v-if="input.api == 'baidu'">
        <span class="label">语速（0-15）</span>
        <el-slider :min="0" :max="15" :step="1" v-model="input.baidu.speed"></el-slider>
        <span class="label">音调（0-15）</span>
        <el-slider :min="0" :max="15" :step="1" v-model="input.baidu.pitch"></el-slider>
        <span class="label">音量（0-15）</span>
        <el-slider :min="0" :max="15" :step="1" v-model="input.baidu.volumn"></el-slider>
        <span class="label">朗读角色：</span>
        <div>
          <el-select v-model="input.baidu.person" style="width: 100%; margin: 5px 0;">
            <el-option label="普通女声" :value="0"></el-option>
            <el-option label="普通男声" :value="1"></el-option>
            <el-option label="合成-度逍遥" :value="3"></el-option>
            <el-option label="合成-度丫丫" :value="4"></el-option>
          </el-select>
        </div>
      </div>
      <div v-if="input.api == 'xfyun'">
        <span class="label">语速（0-100）</span>
        <el-slider :min="0" :max="100" :step="1" v-model="input.xfyun.speed"></el-slider>
        <span class="label">音调（0-100）</span>
        <el-slider :min="0" :max="100" :step="1" v-model="input.xfyun.pitch"></el-slider>
        <span class="label">音量（0-100）</span>
        <el-slider :min="0" :max="100" :step="1" v-model="input.xfyun.volumn"></el-slider>
        <span class="label">朗读角色：</span>
        <div>
          <el-select v-model="input.xfyun.voiceName" style="width: 100%; margin: 5px 0;">
            <el-option label="讯飞小燕" value="xiaoyan"></el-option>
            <el-option label="讯飞许久" value="aisjiuxu"></el-option>
            <el-option label="讯飞小萍" value="aisxping"></el-option>
            <el-option label="讯飞小婧" value="aisjinger"></el-option>
            <el-option label="讯飞许小宝" value="aisbabyxu"></el-option>
          </el-select>
        </div>
        <span class="label">引擎类型：</span>
        <div>
          <el-select v-model="input.xfyun.engineType" style="width: 100%; margin: 5px 0;">
            <el-option label="普通效果" value="aisound"></el-option>
            <el-option label="中文" value="intp65"></el-option>
            <el-option label="英文" value="intp65_en"></el-option>
            <el-option label="小语种" value="mtts"></el-option>
            <el-option label="优化效果" value="x"></el-option>
          </el-select>
        </div>
      </div>
    </div>

    <div class="grid content-uploader">
      <h3>文本上传</h3>
      <span class="label">选择上传方式：</span>
      <el-radio-group v-model="input.contentType">
        <el-radio-button
          v-for="(type, ti) in contentTypes.filter(c => !c.disabled)"
          :key="ti"
          :label="type.value"
        >{{type.label}}</el-radio-button>
      </el-radio-group>
      <div class="content" style="margin-top: 20px;" v-if="input.contentType == 'text'">
        <textarea v-model.trim="input.text" placeholder="输入需要朗读的文本" name="" id="" cols="30" rows="10"></textarea>
      </div>
      <div class="content" style="margin-top: 20px;" v-if="input.contentType == 'file'">
        <el-upload
          class="uploader"
          drag
          v-show="false"
        >
          <i class="el-icon-upload"></i>
          <div class="el-upload__text">将文件拖拽到此处，或<em>点击上传</em></div>
          <div class="el-upload__tip" slot="tip">只能上传txt文件</div>
        </el-upload>
        <el-alert>暂不支持</el-alert>
      </div>
      <div style="margin-top: 20px">
        <el-button class="el-icon-service" @click="genAudioFile"> 生成音频</el-button>
      </div>
      <div style="margin-top: 20px" v-if="audioFileToken">
        <div>
          <audio controls>
            <source :src="`/ai/tts/audio-file/${audioFileToken.tag}`" type="audio/mpeg">
          </audio>
        </div>
        <el-button class="el-icon-download" @click="tryDownloadAudio"> 下载音频</el-button>
      </div>
    </div>
  </div>
</template>

<script>
import {download} from './download.js'

export default {
  data() {
    let apiList = [
      {value: 'xfyun', label: '科大讯飞'},
      {value: 'baidu', label: '百度语音'},
    ]
    let contentTypes = [
      {value: 'text', label: '输入文字', disabled: false},
      {value: 'file', label: '上传文件', disabled: true },
    ]
    return {
      apiList,
      contentTypes,
      audioFileToken: '',
      input: {
        api: apiList[0].value,
        contentType: contentTypes[0].value,
        text: '',
        file: '',
        xfyun: {
          voiceName: 'xiaoyan',
          speed: 50,
          volumn: 50,
          pitch: 50,
          engineType: 'intp65',
        },
        baidu: {
          person: 1,
          speed: 5,
          volumn: 5,
          pitch: 5,
        }
      }
    }
  },
  methods: {
    async genAudioFile() {
      try {
        this.audioFileToken = null
        let res = await fetch('/ai/tts/audio-file', {
          method: 'POST',
          body: JSON.stringify(this.input),
        })
        if (res.status != 200) {
          throw `网络错误（${res.status}）`
        }
        let json = await res.json()
        if (json.errCode != 0) throw json.errInfo

        this.audioFileToken = json.data
        this.$message.success(`音频文件生成成功`)
      } catch (ex) {
        this.$message.error(`生成音频文件失败。\n${ex}`)
      }
    },
    async tryDownloadAudio() {
      try {
        await this.downloadAudio()
      } catch (ex) {
        this.$message.error(`下载失败。\n${ex}`)
      }
    },
    async downloadAudio() {
      let res = await fetch(`/ai/tts/audio-file/${this.audioFileToken.tag}`, {
        method: 'GET',
      })
      if (res.status != 200) {
        let reason = await res.text()
        throw reason
      }

      let blob = await res.blob()
      download(blob, 'audio.mp3')
    }
  }
}
</script>


<style lang="less" scoped>
.grid {
  padding: 10px;
}
span.label {
  font-size: 12px;
}
.t2a-container {
  max-width: 1024px;
  margin: 30px auto;
  text-align: left;
  display: flex;
  flex-direction: row;

  .base-settings {
    
  }
  .content-uploader {
    flex: 1;

    .content {
      textarea {
        width: 100%;
        height: 250px;
        border-radius: 3px;
        border: 1px solid #ddd;
        resize: vertical;
        padding: 1em;
        font-size: 14px;
        line-height: 1.5;
      }
    }
  }
}
</style>
