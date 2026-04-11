/**
 * ComfyUI 图片信息解析工具
 * 解析 ComfyUI 生成的 txt 文件，提取提示词和参数
 */

/**
 * 解析 ComfyUI 图片信息文件
 * @param {string} content - 文件内容
 * @returns {object} 解析结果 { posText, negText, model, params }
 */
export function parseComfyUIFile(content) {
  const lines = content.split('\n').filter(line => line.trim())
  
  if (lines.length < 2) {
    throw new Error('文件格式不正确，需要至少包含 fileinfo 和 prompt 两行')
  }
  
  // 查找 prompt 行（第二行或包含 "prompt"= 的行）
  let promptLine = lines.find(line => line.startsWith('"prompt"='))
  if (!promptLine && lines.length >= 2) {
    promptLine = lines[1] // 默认第二行
  }
  
  if (!promptLine) {
    throw new Error('未找到 prompt 数据')
  }
  
  // 提取 JSON 部分（去掉 "prompt"= 前缀）
  const jsonStr = promptLine.replace(/^"prompt"="/, '').replace(/"$/, '')
  
  // 处理转义的引号
  const unescapedStr = jsonStr.replace(/\\"/g, '"').replace(/\\\\/g, '\\')
  
  let promptData
  try {
    promptData = JSON.parse(unescapedStr)
  } catch (e) {
    // 尝试直接解析（如果没有外层引号）
    try {
      promptData = JSON.parse(promptLine.replace(/^"prompt"=/, ''))
    } catch (e2) {
      throw new Error('解析 prompt JSON 失败: ' + e.message)
    }
  }
  
  return extractPromptInfo(promptData)
}

/**
 * 从 ComfyUI prompt 数据中提取提示词信息
 * @param {object} promptData - 解析后的 prompt JSON
 * @returns {object} 提取的信息 { posText, negText, model, params }
 */
function extractPromptInfo(promptData) {
  const result = {
    posText: '',
    negText: '',
    model: '',
    params: {
      steps: 30,
      cfg: 7,
      sampler: 'DPM++ 2M Karras',
      width: 512,
      height: 512,
    }
  }
  
  const nodes = Object.values(promptData)
  
  // 1. 提取正向和负向提示词（CLIPTextEncode 节点）
  const clipTextNodes = nodes.filter(node => 
    node.class_type === 'CLIPTextEncode' || 
    node.class_type === 'CLIPTextEncodeFlux' ||
    node.class_type?.includes('CLIPText')
  )
  
  // 通常第一个是正向，第二个是负向（根据节点顺序或连接关系）
  if (clipTextNodes.length >= 1) {
    result.posText = clipTextNodes[0].inputs?.text || ''
  }
  if (clipTextNodes.length >= 2) {
    result.negText = clipTextNodes[1].inputs?.text || ''
  }
  
  // 2. 提取模型信息（CheckpointLoaderSimple 节点）
  const checkpointNodes = nodes.filter(node =>
    node.class_type === 'CheckpointLoaderSimple' ||
    node.class_type === 'CheckpointLoader' ||
    node.class_type?.includes('Checkpoint')
  )
  
  if (checkpointNodes.length > 0) {
    const ckptName = checkpointNodes[0].inputs?.ckpt_name || ''
    // 提取文件名（去掉路径）
    result.model = ckptName.split('\\').pop().split('/').pop()
  }
  
  // 3. 提取采样器参数（KSampler 节点）
  const samplerNodes = nodes.filter(node =>
    node.class_type === 'KSampler' ||
    node.class_type === 'KSamplerAdvanced' ||
    node.class_type?.includes('Sampler')
  )
  
  if (samplerNodes.length > 0) {
    const inputs = samplerNodes[0].inputs || {}
    result.params.steps = inputs.steps || 30
    result.params.cfg = inputs.cfg || 7
    result.params.sampler = inputs.sampler_name || inputs.sampler || 'DPM++ 2M Karras'
  }
  
  // 4. 提取图片尺寸（EmptyLatentImage 节点）
  const latentNodes = nodes.filter(node =>
    node.class_type === 'EmptyLatentImage' ||
    node.class_type?.includes('EmptyLatent')
  )
  
  if (latentNodes.length > 0) {
    const inputs = latentNodes[0].inputs || {}
    result.params.width = inputs.width || 512
    result.params.height = inputs.height || 512
  }
  
  // 清理提示词中的多余空格和换行
  result.posText = cleanPromptText(result.posText)
  result.negText = cleanPromptText(result.negText)
  
  return result
}

/**
 * 清理提示词文本
 * @param {string} text - 原始提示词
 * @returns {string} 清理后的提示词
 */
function cleanPromptText(text) {
  if (!text) return ''
  
  return text
    .replace(/\n+/g, ' ')      // 将多个换行替换为空格
    .replace(/\s+/g, ' ')       // 将多个空格合并为一个
    .replace(/,\s*,/g, ',')     // 清理多余的逗号
    .trim()
}

/**
 * 读取文件内容
 * @param {File} file - 文件对象
 * @returns {Promise<string>} 文件内容
 */
export function readFileContent(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.onerror = (e) => reject(new Error('读取文件失败'))
    reader.readAsText(file)
  })
}
