一个剧集推流工具
==========================

- 使用场景
  - 不想使用GUI推流工具
  - 剧集的自动循环推流
  - 自定义背景，自定义正在播放的名称显示

- 使用方法
  - 下载对应平台可执行文件，在其目录下建立`config.json`，配置参数后运行
    
    ```
    ./loop-stream
    ```

  - 配置参数说明

    <table>
      <thead>
        <th>参数</th>
        <th>说明</th>
      </thread>
      <tr><td>ffmpeg</td><td>ffmpeg 可执行程序路径</td></tr>
      <tr><td>ffprobe</td><td>ffprobe 可执行程序路径</td></tr>
      <tr><td>stage.image</td><td>场景的背景图片路径</td></tr>
      <tr><td>stage.width</td><td>场景的宽度</td></tr>
      <tr><td>stage.height</td><td>场景的高度</td></tr>
      <tr><td>input.rectangle</td><td>视频的位置信息, [x, y, width, height]</td></tr>
      <tr><td>input.episodes</td><td>视频剧集文件列表</td></tr>
      <tr><td>input.title.font</td><td>当前播放字体文件路径</td></tr>
      <tr><td>input.title.prefix</td><td>当前播放视频名称显示前缀</td></tr>
      <tr><td>input.title.x</td><td>当前播放显示位置x坐标，支持ffmpeg参数写法</td></tr>
      <tr><td>input.title.y</td><td>当前播放显示位置y坐标，支持ffmpeg参数写法</td></tr>
      <tr><td>input.title.color</td><td>当前播放字体颜色</td></tr>
      <tr><td>input.title.size</td><td>当前播放字体大小</td></tr>
      <tr><td>output.stream_url</td><td>推流地址</td></tr>
    </table>

- 演示效果

  <img width="644" alt="image" src="https://user-images.githubusercontent.com/1967804/166853113-1cb60766-5907-43ac-bf6c-6c59f29c0a0d.png">
