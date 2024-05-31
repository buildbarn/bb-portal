import dayjs from 'dayjs';
import AdvancedFormat from 'dayjs/plugin/advancedFormat';
import Duration from 'dayjs/plugin/duration';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import RelativeTime from 'dayjs/plugin/relativeTime';
import Timezone from 'dayjs/plugin/timezone';
import customParseFormat from 'dayjs/plugin/customParseFormat';
import localeData from 'dayjs/plugin/localeData';
import weekday from 'dayjs/plugin/weekday';
import weekOfYear from 'dayjs/plugin/weekOfYear';
import weekYear from 'dayjs/plugin/weekYear';

dayjs.extend(customParseFormat);
dayjs.extend(weekday);
dayjs.extend(localeData);
dayjs.extend(weekOfYear);
dayjs.extend(weekYear);
dayjs.extend(AdvancedFormat);
dayjs.extend(Duration);
dayjs.extend(LocalizedFormat);
dayjs.extend(RelativeTime);
dayjs.extend(Timezone);

export default dayjs;
