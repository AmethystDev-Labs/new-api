/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import zhCNTranslation from './locales/zh-CN.json';
const DEFAULT_LANGUAGE = 'zh-CN';
const localeLoaders = import.meta.glob('./locales/*.json');
import { normalizeLanguage, supportedLanguages } from './language';

async function ensureLanguageResources(language) {
  const normalizedLanguage = normalizeLanguage(language);
  if (i18n.hasResourceBundle(normalizedLanguage, 'translation')) {
    return normalizedLanguage;
  }

  const loader = localeLoaders[`./locales/${normalizedLanguage}.json`];
  if (!loader) {
    return DEFAULT_LANGUAGE;
  }

  try {
    const module = await loader();
    const translation =
      module?.default?.translation || module?.translation || module?.default;
    if (translation) {
      i18n.addResourceBundle(
        normalizedLanguage,
        'translation',
        translation,
        true,
        true,
      );
    }
  } catch (error) {
    console.error('Failed to load locale:', normalizedLanguage, error);
    return DEFAULT_LANGUAGE;
  }

  return normalizedLanguage;
}

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    load: 'currentOnly',
    supportedLngs: supportedLanguages,
    resources: {
      [DEFAULT_LANGUAGE]: zhCNTranslation,
    },
    fallbackLng: DEFAULT_LANGUAGE,
    nsSeparator: false,
    interpolation: {
      escapeValue: false,
    },
  });

const initialLanguage = normalizeLanguage(i18n.resolvedLanguage || i18n.language);
void ensureLanguageResources(initialLanguage).then((loadedLanguage) => {
  if (loadedLanguage !== i18n.language) {
    i18n.changeLanguage(loadedLanguage);
  }
});

i18n.on('languageChanged', (language) => {
  void ensureLanguageResources(language);
});

export default i18n;
