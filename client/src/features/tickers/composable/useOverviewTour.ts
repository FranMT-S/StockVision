import { ChartType } from '@/shared/enums/chart';
import Shepherd, { Step, StepOptions } from 'shepherd.js';
import 'shepherd.js/dist/css/shepherd.css';
import { Ref, nextTick } from 'vue';
import { useOnboardingStore } from '@/shared/store/onboardingStore';
import { storeToRefs } from 'pinia';

export const useOverviewTour = (chartTypeRef: Ref<ChartType>,options?:{
  cancelIcon?: boolean,
}) => {
  const {updateOnboarding} = useOnboardingStore()
  const {overviewData} = storeToRefs(useOnboardingStore())
  const chartType = chartTypeRef;
  const {cancelIcon = true} = options ?? {}

  const tour = new Shepherd.Tour({
    useModalOverlay: true,
    defaultStepOptions: {
      cancelIcon: { enabled: cancelIcon },
      canClickTarget: false,
      scrollTo: { behavior: 'smooth', block: 'center' },
      classes: 'custom-tour' 
    }
  });

  // update the cancel icon for all steps
  function setCancelIcon(enabled: boolean) {
    tour.steps.forEach(step => {
      step.options.cancelIcon = { enabled };
    });
  }

  /** container to show friendly indicator of steps  */
  const CreateBreadcrumbSteps = () =>{
      const container = tour.currentStep?.getElement()?.querySelector('.shepherd-content');
      if (!container) return null;

      let indicatorContainer = container.querySelector('.tw-step-indicators');
      if (indicatorContainer) indicatorContainer.remove();

      indicatorContainer = document.createElement('div');
      indicatorContainer.className = 'tw-step-indicators tw-flex tw-space-x-1 tw-mt-2 tw-justify-center';

      const currentIndex = tour.steps.indexOf(tour.currentStep!);

      for (let i = 0; i < tour.steps.length; i++) {
        const circle = document.createElement('span');
        circle.className = `
          tw-w-2 tw-h-2 tw-rounded-full
          ${i === currentIndex ? 'tw-bg-blue-500' : 'tw-bg-gray-300'}
        `;
        indicatorContainer.appendChild(circle);
      }

      return indicatorContainer;
    }

  // animate popover, update smooth transition when shown
  const animatePopover = (step: Step) => {
    const el = step.getElement();
    if (!el) return;

    el.style.opacity = '0';
    el.style.transform = 'translateY(8px)';
    el.style.transition = 'opacity 0.5s ease, transform 0.5s ease';

    void el.offsetHeight;

    nextTick(() => {
      el.style.opacity = '1';
      el.style.transform = 'translateY(0)';
    });

    const overlay = document.querySelector('.shepherd-modal-overlay-container') as HTMLElement;
    if (overlay) {
      overlay.style.opacity = '0';
      overlay.style.transition = 'opacity 0.5s ease';
      void overlay.offsetHeight;
      nextTick(() => {
        overlay.style.opacity = '0.5';
      });
    }
  };

  /** add step to tour  with animate popover and breadcrumb steps */
  const addStep = (options: StepOptions | Step) => {
    tour.addStep({
      ...options,
      when: {
        show: () => {
          const contentContainer = tour.currentStep?.getElement()?.querySelector('.shepherd-content');
          if (!contentContainer) 
            return;

          animatePopover(tour.currentStep!);
          
          const indicatorContainer = CreateBreadcrumbSteps();
          if (!indicatorContainer) 
            return;
          
          contentContainer.appendChild(indicatorContainer);
        }
      }
    });
  };

  /** action when next button is pressed, update onboarding data */
  const nextStep = async () => {
    const index = tour.steps.indexOf(tour.currentStep!);
    if(!overviewData.value?.overviewDone){
      try {
        const isFinish = index === tour.steps.length - 2;
        await updateOnboarding(index + 2, isFinish)
      } catch (error) {
        console.error(error)
      }
    }
    tour.next()
  }

  tour.on('complete', async () => {
    overviewData.value= {id: 1, overviewStep: tour.steps.length, overviewDone: true };
    await updateOnboarding(tour.steps.length, true)
  },undefined, true);


  // Step 1 - Company Info
  addStep({
    id: 'company-info',
    title: 'Company Information',
    text: 'View the company’s details, key metrics, and today’s recommendations. You can also visit the company website from here.',
    attachTo: { element: '#company-info-col', on: 'right-start' },
    buttons: [{ text: 'Next', action:async () => await nextStep() }]
  });

  // Step 2 - Charts
  addStep({
    id: 'charts',
    title: 'Charts',
    text: 'This chart displays the stock information. You can adjust the time interval and take a screenshot of the chart.',
    attachTo: { element: '.chart-container', on: 'left' },
    buttons: [{ text: 'Next', action:async () => await nextStep() }],
    extraHighlights: ['.controls-right', '.controls-left']
  });

  // Step 3 - Candlestick
  addStep({
    id: 'candlestick',
    title: 'Candlestick Chart',
    text: 'The candlestick chart shows the stock’s price history and trading volume at the bottom.',
    attachTo: { element: '#candlestick-btn', on: 'left' },
    extraHighlights: ['.chart-container'],
    buttons: [
      { text: 'Next', action:async () => { chartType.value = ChartType.area; await nextStep(); } }
    ]
  });

  // Step 4 - Area
  addStep({
    id: 'close-price',
    title: 'Close Price Chart',
    text: 'The area chart displays the stock’s closing prices over time.',
    attachTo: { element: '#area-btn', on: 'left' },
    extraHighlights: ['.chart-container'],
    buttons: [
      { text: 'Next', action:async () => { chartType.value = ChartType.candlestick; await nextStep(); } }
    ]
  });

  // Step 5 - Predict with IA
  addStep({
    id: 'predict',
    title: 'Predict with IA',
    text: 'See the AI prediction of the stock’s future performance for the upcoming week.',
    attachTo: { element: '#predict-btn', on: 'right' },
    buttons: [
      { text: 'Next', action:async () => await nextStep() }
    ]
  });

  // Step 6 - Recommendations
  addStep({
    id: 'recommendations',
    title: 'Recommendations',
    text: 'See broker recommendations on stocks to guide your investment decisions.',
    attachTo: { element: '.recommendations-table', on: 'top' },
    buttons: [{ text: 'Next', action:async () => await nextStep() }]
  });

  // Step 7 - Company News
  addStep({
    id: 'company-news',
    title: 'Company News',
    text: 'Check the company news from the past year to stay updated.',
    attachTo: { element: '.company-news-sidebar-title', on: 'left-start' },
    buttons: [{ text: 'Finish', action: tour.complete }],
    extraHighlights: ['.company-news-sidebar']
  });

  return {
    overviewTour: tour,
    setCancelIcon
  };
};
