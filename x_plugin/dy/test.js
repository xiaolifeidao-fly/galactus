

function simulateTyping(inputElement, text, interval = 100) {
  let index = 0;

  // 清空输入框（可选）
  inputElement.value = '';

  // 定时器模拟逐个输入
  const typingInterval = setInterval(() => {
    // 添加字符
    inputElement.value += text[index];
    const inputEvent = new Event('input', {
        isTrusted : true,
        bubbles: true,      // 事件是否冒泡
        cancelBubble : false,
        cancelable: false,   // 事件是否可取消
        composed : true,
        data: inputElement.value,  // 输入框的新值
        inputType: 'insertText',
        isComposing: false,
        defaultPrevented : false,
        detail : 0,
        eventPhase : 3,
        isComposing : false,
        returnValue : true,
        timeStamp : 1020156.799999997,
        which  : 0
      });
    inputElement.dispatchEvent(inputEvent);
    // 如果已经输入完所有字符，清除定时器
    index++;
    if (index === text.length) {
      clearInterval(typingInterval);
    }
  }, interval); // 控制输入速度，单位为毫秒
}

// 使用示例
const phoneInput = document.getElementsByClassName("web-login-normal-input__input")[0];
simulateTyping(phoneInput, '18217637991', 200); 


