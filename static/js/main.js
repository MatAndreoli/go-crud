/**
 * main.js
 * Lógica de frontend aprimorada com IMask.js, validação real-time e modal.
 */

document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("crud-form");
  if (!form) return;

  const formTitle = document.getElementById("form-title");
  const formSubmitBtn = document.getElementById("form-submit-btn");
  const formCancelBtn = document.getElementById("form-cancel-btn");
  const formIdField = document.getElementById("form-id-field");
  const formInputs = form.querySelectorAll("input[name]");

  // Elementos da modal
  const modalOverlay = document.getElementById("modal-overlay");
  const openModalBtn = document.getElementById("open-modal-btn");
  const closeModalBtn = document.getElementById("close-modal-btn");

  const originalFormAction = form.action;
  const originalFormTitle = formTitle.innerText;
  const originalSubmitText = formSubmitBtn.innerText;
  let inputMasks = [];

  // Funções da Modal
  const openModal = () => {
    modalOverlay.classList.remove("hidden");
    document.body.style.overflow = "hidden";
  };

  const closeModal = () => {
    modalOverlay.classList.add("hidden");
    document.body.style.overflow = "";
  };

  // Event Listeners da Modal
  openModalBtn.addEventListener("click", () => {
    cancelEdit(); // Limpa o formulário
    openModal();
  });

  closeModalBtn.addEventListener("click", closeModal);

  formCancelBtn.addEventListener("click", () => {
    cancelEdit();
    closeModal();
  });

  // Fechar modal ao clicar no overlay
  modalOverlay.addEventListener("click", (e) => {
    if (e.target === modalOverlay) {
      cancelEdit();
      closeModal();
    }
  });

  // Fechar modal com ESC
  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape" && !modalOverlay.classList.contains("hidden")) {
      cancelEdit();
      closeModal();
    }
  });

  const initMasks = () => {
    const customDefinitions = {
      9: {
        mask: /\d/,
        lazy: false,
      },

      "#": {
        mask: /[a-zA-Z]/,
        lazy: false,
      },

      "*": {
        mask: /./,
        lazy: false,
      },
    };
    inputMasks = [];
    formInputs.forEach((input) => {
      const maskPattern = input.dataset.mask;
      if (!maskPattern) return;

      const imaskPattern = maskPattern.replace(/9/g, "0");

      let maskOptions = {
        mask: imaskPattern,
        lazy: false,
        definitions: customDefinitions,
      };

      if (maskPattern.includes("(99) 9")) {
        maskOptions.mask = [
          {
            mask: "(00) 0000-0000",
            lazy: false,
          },
          {
            mask: "(00) 00000-0000",
            lazy: false,
          },
        ];
      }

      const maskInstance = IMask(input, maskOptions);
      inputMasks.push(maskInstance);
    });
  };

  const initValidation = () => {
    formInputs.forEach((input) => {
      input.addEventListener("blur", (e) => {
        validateField(e.target);
      });

      input.addEventListener("input", (e) => {
        clearError(e.target);
      });
    });
  };

  const validateField = (input) => {
    const value = input.value;
    const type = input.dataset.validateType;
    const isRequired = input.hasAttribute("required");

    if (isRequired && value.trim() === "") {
      showError(input, "Campo obrigatório");
      return false;
    }

    if (value.trim() !== "") {
      switch (type) {
        case "cpf":
          if (!isValidCPF(value, isRequired)) {
            showError(input, "CPF inválido");
            return false;
          }
          break;
        case "email":
          if (!isValidEmail(value, isRequired)) {
            showError(input, "Email inválido");
            return false;
          }
          break;
        case "cnpj":
          if (!isValidCNPJ(value, isRequired)) {
            showError(input, "CNPJ inválido");
            return false;
          }
          break;
        case "telefone":
          if (!isValidTelefone(value, isRequired)) {
            showError(input, "Telefone inválido");
            return false;
          }
          break;
        case "cep":
          if (!isValidCEP(value, isRequired)) {
            showError(input, "CEP inválido");
            return false;
          }
          break;
      }
    }

    clearError(input);
    return true;
  };

  const showError = (input, message) => {
    const errorJS = document.getElementById(`error-js-${input.name}`);
    const errorBackend = document.getElementById(`error-backend-${input.name}`);

    if (errorBackend) errorBackend.classList.add("hidden");

    if (errorJS) {
      errorJS.innerText = message;
      errorJS.classList.remove("hidden");
    }
    input.classList.add("border-red-500", "ring-1", "ring-red-500");
  };

  const clearError = (input) => {
    const errorJS = document.getElementById(`error-js-${input.name}`);
    const errorBackend = document.getElementById(`error-backend-${input.name}`);

    if (errorBackend) errorBackend.classList.add("hidden");
    if (errorJS) errorJS.classList.add("hidden");

    input.classList.remove("border-red-500", "ring-1", "ring-red-500");
  };

  const clearAllValidation = () => {
    formInputs.forEach((input) => clearError(input));
  };

  const validateForm = () => {
    let isFormValid = true;
    formInputs.forEach((input) => {
      if (!validateField(input)) {
        isFormValid = false;
      }
    });
    return isFormValid;
  };

  window.startEdit = async (id) => {
    try {
      const response = await fetch(`/get?id=${id}`);
      if (!response.ok) throw new Error("Falha ao carregar dados");

      const data = await response.json();

      cancelEdit();

      formInputs.forEach((input) => {
        if (data[input.name]) {
          let value = data[input.name];

          if (input.type === "date" && value) {
            value = value.split("T")[0];
          }

          input.value = value;

          const mask = inputMasks.find((m) => m.el === input);
          if (mask) {
            mask.updateValue();
          }
        }
      });

      formIdField.value = id;
      form.action = `/update?id=${id}`;
      formTitle.innerText = `Editando Registro #${id}`;
      formSubmitBtn.innerText = "Atualizar";

      openModal();
    } catch (error) {
      console.error("Falha ao buscar dados para edição:", error);
      alert("Não foi possível carregar os dados para edição.");
    }
  };

  const cancelEdit = () => {
    form.reset();
    clearAllValidation();

    inputMasks.forEach((mask) => mask.updateValue());

    formIdField.value = "";
    form.action = originalFormAction;
    formTitle.innerText = originalFormTitle;
    formSubmitBtn.innerText = originalSubmitText;
  };

  form.addEventListener("submit", (e) => {
    if (!validateForm()) {
      e.preventDefault();
      console.warn("Formulário inválido. Verifique os campos.");

      form.querySelector(".is-invalid")?.focus();
    }
  });

  const isValidEmail = (email, isRequired) => {
    if (email === "" && !isRequired) return true;

    return /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(email);
  };

  const isValidCPF = (cpf, isRequired) => {
    cpf = cpf.replace(/\D/g, "");

    if (cpf === "" && !isRequired) return true;

    if (cpf.length !== 11 || /^(\d)\1{10}$/.test(cpf)) return false;

    let sum = 0,
      rest;

    for (let i = 1; i <= 9; i++)
      sum += parseInt(cpf.substring(i - 1, i)) * (11 - i);
    rest = (sum * 10) % 11;
    if (rest === 10 || rest === 11) rest = 0;
    if (rest !== parseInt(cpf.substring(9, 10))) return false;

    sum = 0;
    for (let i = 1; i <= 10; i++)
      sum += parseInt(cpf.substring(i - 1, i)) * (12 - i);
    rest = (sum * 10) % 11;
    if (rest === 10 || rest === 11) rest = 0;
    if (rest !== parseInt(cpf.substring(10, 11))) return false;

    return true;
  };

  const isValidCNPJ = (cnpj, isRequired) => {
    cnpj = cnpj.replace(/[^\d]+/g, "");

    if (cnpj === "" && !isRequired) return true;

    if (cnpj.length !== 14) return false;

    if (/^(.)\1{13}$/.test(cnpj)) return false;

    let soma = 0;
    let peso = [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2];
    for (let i = 0; i < 12; i++) {
      soma += parseInt(cnpj[i]) * peso[i + 1];
    }
    let digito1 = 11 - (soma % 11);
    digito1 = digito1 === 10 || digito1 === 11 ? 0 : digito1;

    soma = 0;
    peso = [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2];
    for (let i = 0; i < 13; i++) {
      soma += parseInt(cnpj[i]) * peso[i];
    }
    let digito2 = 11 - (soma % 11);
    digito2 = digito2 === 10 || digito2 === 11 ? 0 : digito2;

    return cnpj[12] == digito1 && cnpj[13] == digito2;
  };

  const isValidTelefone = (telefone, isRequired) => {
    telefone = telefone.replace(/\D/g, "");

    if (telefone === "" && !isRequired) return true;

    return telefone.length >= 10 && telefone.length <= 11;
  };

  const isValidCEP = (cep, isRequired) => {
    cep = cep.replace(/\D/g, "");

    if (cep === "" && !isRequired) return true;

    return cep.length === 8;
  };

  initMasks();
  initValidation();

  // Abrir modal automaticamente se houver erros de backend
  const hasBackendErrors = document.querySelector('[id^="error-backend-"]');
  if (hasBackendErrors && !hasBackendErrors.classList.contains("hidden")) {
    openModal();
  }
});